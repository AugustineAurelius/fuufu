package cmd

import (
	"net"
	"net/http"

	"github.com/AugustineAurelius/fuufu/api/todo"
	"github.com/AugustineAurelius/fuufu/frontend"
	"github.com/AugustineAurelius/fuufu/internal/config"
	todo_repository "github.com/AugustineAurelius/fuufu/internal/repository/todo"
	"github.com/AugustineAurelius/fuufu/internal/server"
	"github.com/AugustineAurelius/fuufu/pkg/common"
	"github.com/AugustineAurelius/fuufu/pkg/middleware"
	"github.com/spf13/cobra"
)

func createServeCMD(manager *config.Manager) *cobra.Command {
	var debug, dev bool

	serveCMD := &cobra.Command{
		Use: "serve",
		Run: func(cmd *cobra.Command, args []string) {
			log := getLogger(manager)

			postgresMasterConfig := manager.LoadPostgres()
			log.Sugar().Infof("get postgres connection url: %s", postgresMasterConfig.URL())

			postgresSlaveConfig := manager.LoadPostgresSlave()
			log.Sugar().Infof("get postgres connection url: %s", postgresSlaveConfig.URL())

			pgMaster, err := common.NewPostgres(cmd.Context(), postgresMasterConfig, log)
			if err != nil {
				log.Panic(err.Error())
			}
			if err = pgMaster.Pool.Ping(cmd.Context()); err != nil {
				log.Panic(err.Error())
			}

			pgSlave, err := common.NewPostgres(cmd.Context(), postgresSlaveConfig, log)
			if err != nil {
				log.Panic(err.Error())
			}
			if err = pgSlave.Pool.Ping(cmd.Context()); err != nil {
				log.Panic(err.Error())
			}

			todoRepository := todo_repository.New(&pgMaster)

			todoHandlers := todo.NewStrictHandler(&server.TodoHandler{
				Repo: todoRepository,
			}, nil)

			r := http.NewServeMux()
			h := todo.HandlerFromMux(todoHandlers, r)

			h = middleware.LoggingMiddleware(log, h)

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
