package cmd

import (
	"net/http"

	"github.com/AugustineAurelius/fuufu/api/auth"
	"github.com/AugustineAurelius/fuufu/api/todo"
	"github.com/AugustineAurelius/fuufu/frontend"
	"github.com/AugustineAurelius/fuufu/internal/analytic"
	"github.com/AugustineAurelius/fuufu/internal/config"
	event_repository "github.com/AugustineAurelius/fuufu/internal/repository/event"
	todo_repository "github.com/AugustineAurelius/fuufu/internal/repository/todo"
	user_repository "github.com/AugustineAurelius/fuufu/internal/repository/user"

	"github.com/AugustineAurelius/fuufu/internal/server"
	"github.com/AugustineAurelius/fuufu/pkg/common"
	"github.com/AugustineAurelius/fuufu/pkg/geo"
	"github.com/AugustineAurelius/fuufu/pkg/logger"
	"github.com/AugustineAurelius/fuufu/pkg/middleware"
	"github.com/AugustineAurelius/fuufu/pkg/migration"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.14.0"
)

func createServeCMD(manager *config.Manager) *cobra.Command {
	var debug, dev bool

	serviceName := semconv.ServiceNameKey.String("FUUFU")
	name := "fuufu"
	serveCMD := &cobra.Command{
		Use: "serve",
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.NewWithManager(manager)

			if err := migration.CheckMigrations(cmd.Context(), manager.LoadPostgres()); err != nil {
				log.Error(err.Error())
				return
			}

			conn, err := initCollector(manager)
			if err != nil {
				log.Error(err.Error())
				return
			}

			res, err := resource.New(cmd.Context(), resource.WithAttributes(serviceName))
			if err != nil {
				log.Error(err.Error())
				return
			}

			shutdownTracerProvider, err := initTracerProvider(cmd.Context(), res, conn)
			if err != nil {
				log.Error(err.Error())
				return
			}

			shutdownMetricProvider, err := initMeterProvider(cmd.Context(), res, conn)
			if err != nil {
				log.Error(err.Error())
				return
			}
			defer func() {
				shutdownTracerProvider(cmd.Context())
				shutdownMetricProvider(cmd.Context())
			}()

			tracer := otel.Tracer(name)
			meter := otel.Meter(name)

			exporter, err := prometheus.New()
			if err != nil {
				panic(err)
			}

			// Set the global MeterProvider to the Prometheus exporter
			meterProvider := metric.NewMeterProvider(
				metric.WithReader(exporter),
			)
			otel.SetMeterProvider(meterProvider)
			err = runtime.Start()
			if err != nil {
				log.Error(err.Error())
				return
			}
			postgresMasterConfig := manager.LoadPostgres()
			postgresSlaveConfig := manager.LoadPostgresSlave()

			pgMaster, err := common.NewPostgres(cmd.Context(), postgresMasterConfig, log, tracer)
			if err != nil {
				log.Error(err.Error())
				return
			}
			if err = pgMaster.Pool.Ping(cmd.Context()); err != nil {
				log.Error(err.Error())
				return
			}
			log.Info("successfully conected to postgresMaster")

			pgSlave, err := common.NewPostgres(cmd.Context(), postgresSlaveConfig, log, tracer)
			if err != nil {
				log.Error(err.Error())
				return
			}
			if err = pgSlave.Pool.Ping(cmd.Context()); err != nil {
				log.Error(err.Error())
				return
			}
			log.Info("successfully conected to postgresSlave")

			todoRepository := todo_repository.NewCommand(&pgMaster)
			userRepository := user_repository.NewCommand(&pgMaster)
			eventRepository := event_repository.NewCommand(&pgMaster)

			todoQuery := todo_repository.NewQuery(&pgSlave)

			analitic := analytic.Analytic{
				Meter:    meter,
				TodoRepo: todoQuery,
			}

			go func() {
				err = analitic.Run(cmd.Context())
				if err != nil {
					log.Error(err.Error())
					return
				}
			}()

			todoHandlers := todo.NewStrictHandler(&server.TodoHandler{
				Repo:      todoRepository,
				Telemetry: tracer,
			}, nil)

			authHadnlers := auth.NewStrictHandler(&server.AuthHandler{
				Repo:      userRepository,
				Telemetry: tracer,
				Shield:    manager.LoadShield(),
			}, nil)

			r := http.NewServeMux()
			r.Handle("/metrics", promhttp.Handler())
			auth.HandlerFromMux(authHadnlers, r)
			frontend.RegisterFrontend(r)
			h := todo.HandlerFromMux(todoHandlers, r)
			h = middleware.AuthMiddleware(manager.LoadShield(), h, manager.LoadServer().AuthMiddlewareExclude...)

			h = middleware.LoggingMiddleware(log, h)
			h = middleware.MetricMiddleware(meter, h)
			h = middleware.TracingMiddleware(tracer, h)
			h = middleware.AuditMiddleware(geo.NewGeoIPService(), eventRepository, h, manager.LoadServer().GeoMiddlewareExclude...)
			s := &http.Server{
				Handler: h,
				Addr:    manager.LoadServer().Addr,
			}

			panic(s.ListenAndServe())
		},
	}

	serveCMD.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug endpoints")
	serveCMD.PersistentFlags().BoolVar(&dev, "dev", false, "Enable developer options")

	return serveCMD
}
