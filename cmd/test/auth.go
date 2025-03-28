package test

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/AugustineAurelius/fuufu/api/auth"
	"github.com/AugustineAurelius/fuufu/api/todo"
	"github.com/AugustineAurelius/fuufu/internal/config"
	"github.com/AugustineAurelius/fuufu/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func createAuthCMD(manager *config.Manager) *cobra.Command {
	client, err := auth.NewClient("http://" + manager.LoadServer().Addr)
	if err != nil {
		panic(err)
	}
	testCMD := &cobra.Command{
		Use: "auth",
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.NewWithManager(manager)

			resp, err := client.PostApiV1AuthSignup(cmd.Context(), auth.PostApiV1AuthSignupJSONRequestBody{
				Email:    "jon_doe@mail.com",
				Password: "1337",
				Username: "jon doe",
			})
			if err != nil {
				panic(err)
			}

			resp, err = client.PostApiV1AuthSignin(cmd.Context(), auth.PostApiV1AuthSigninJSONRequestBody{
				Password: "1337",
				Username: "jon doe",
			})
			if err != nil {
				panic(err)
			}

			var token auth.PostApiV1AuthSignin200JSONResponse
			if err = json.NewDecoder(resp.Body).Decode(&token); err != nil {
				panic(err)
			}
			log.Info("signin", zap.String("token", token.Token))

			todoClient, err := todo.NewClient("http://" + manager.LoadServer().Addr)
			if err != nil {
				panic(err)
			}

			resp, err = todoClient.GetAllTodos(cmd.Context(), func(ctx context.Context, req *http.Request) error {
				req.Header.Add("Authorization", token.Token)
				return nil
			})
			var todos todo.GetAllTodos200JSONResponse
			if err = json.NewDecoder(resp.Body).Decode(&todos); err != nil {

				panic(err)
			}
			log.Info("get todo", zap.Any("todos", todos.Tasks))
		}}
	return testCMD
}
