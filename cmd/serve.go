package cmd

import (
	"net"
	"net/http"

	"github.com/AugustineAurelius/fuufu/api/todo"
	"github.com/AugustineAurelius/fuufu/frontend"
	"github.com/AugustineAurelius/fuufu/internal/analytic"
	"github.com/AugustineAurelius/fuufu/internal/config"
	todo_repository "github.com/AugustineAurelius/fuufu/internal/repository/todo"
	"github.com/AugustineAurelius/fuufu/internal/server"
	"github.com/AugustineAurelius/fuufu/pkg/common"
	"github.com/AugustineAurelius/fuufu/pkg/middleware"
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
			log := getLogger(manager)

			conn, err := initCollector(manager)
			if err != nil {
				log.Panic(err.Error())
			}

			res, err := resource.New(cmd.Context(), resource.WithAttributes(serviceName))
			if err != nil {
				log.Panic(err.Error())
			}

			shutdownTracerProvider, err := initTracerProvider(cmd.Context(), res, conn)
			if err != nil {
				log.Panic(err.Error())
			}

			shutdownMetricProvider, err := initMeterProvider(cmd.Context(), res, conn)
			if err != nil {
				log.Panic(err.Error())
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
				log.Panic(err.Error())
			}
			postgresMasterConfig := manager.LoadPostgres()
			log.Sugar().Infof("get postgres connection url: %s", postgresMasterConfig.URL())

			postgresSlaveConfig := manager.LoadPostgresSlave()
			log.Sugar().Infof("get postgres connection url: %s", postgresSlaveConfig.URL())

			pgMaster, err := common.NewPostgres(cmd.Context(), postgresMasterConfig, log, tracer)
			if err != nil {
				log.Panic(err.Error())
			}
			if err = pgMaster.Pool.Ping(cmd.Context()); err != nil {
				log.Panic(err.Error())
			}

			pgSlave, err := common.NewPostgres(cmd.Context(), postgresSlaveConfig, log, tracer)
			if err != nil {
				log.Panic(err.Error())
			}
			if err = pgSlave.Pool.Ping(cmd.Context()); err != nil {
				log.Panic(err.Error())
			}

			todoRepository := todo_repository.New(&pgMaster)

			analitic := analytic.Analytic{
				Meter:    meter,
				TodoRepo: todoRepository,
			}

			go func() {
				err = analitic.Run(cmd.Context())
				if err != nil {
					log.Panic(err.Error())
				}
			}()

			todoHandlers := todo.NewStrictHandler(&server.TodoHandler{
				Repo:      todoRepository,
				Telemetry: tracer,
			}, nil)

			r := http.NewServeMux()
			h := todo.HandlerFromMux(todoHandlers, r)

			h = middleware.TracingMiddleware(tracer, h)
			h = middleware.MetricMiddleware(meter, h)
			h = middleware.LoggingMiddleware(log, h)
			r.Handle("/metrics", promhttp.Handler())

			frontend.RegisterFrontend(r)

			s := &http.Server{
				Handler: h,
				Addr:    net.JoinHostPort("0.0.0.0", "7070"),
			}

			panic(s.ListenAndServe())
		},
	}

	serveCMD.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug endpoints")
	serveCMD.PersistentFlags().BoolVar(&dev, "dev", false, "Enable developer options")

	return serveCMD
}
