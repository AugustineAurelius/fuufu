package test

import (
	"github.com/AugustineAurelius/fuufu/api/auth"
	"github.com/AugustineAurelius/fuufu/internal/config"
	"github.com/AugustineAurelius/fuufu/pkg/logger"
	"github.com/oapi-codegen/runtime/types"
	"github.com/spf13/cobra"
)

func createAuthCMD(manager *config.Manager) *cobra.Command {
	testCMD := &cobra.Command{
		Use: "auth",
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.NewWithManager(manager)

			resp, err := authClient.PostApiV1AuthSignupWithResponse(cmd.Context(), auth.PostApiV1AuthSignupJSONRequestBody{
				Email:    types.Email(myEmail),
				Password: myPassword,
				Username: myUsername,
			})
			checkError(err)
			checkResponse(resp.HTTPResponse)

			generateToken(cmd.Context())

			log.Info("success test")
		}}
	return testCMD
}
