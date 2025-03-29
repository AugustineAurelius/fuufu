package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AugustineAurelius/fuufu/api/auth"
	"github.com/AugustineAurelius/fuufu/api/todo"
	"github.com/AugustineAurelius/fuufu/internal/config"
)

var (
	authClient *auth.ClientWithResponses
	todoClient *todo.ClientWithResponses

	myEmail    = "jon_doe@mail.com"
	myPassword = "1337"
	myUsername = "jon doe"
)

func createClients(manager *config.Manager) {
	serverURL := "http://" + manager.LoadServer().Addr

	var err error

	authClient, err = auth.NewClientWithResponses(serverURL)
	checkError(err)

	todoClient, err = todo.NewClientWithResponses(serverURL)
	checkError(err)
}

func checkResponse(resp *http.Response) error {
	if resp.StatusCode != 200 || resp.StatusCode != 409 {
		return fmt.Errorf("unexpected status code")
	}
	return nil
}

func generateToken(ctx context.Context) string {
	resp, err := authClient.PostApiV1AuthSigninWithResponse(ctx, auth.PostApiV1AuthSigninJSONRequestBody{
		Password: myPassword,
		Username: myUsername,
	})
	checkError(err)
	checkResponse(resp.HTTPResponse)

	var token auth.PostApiV1AuthSignin200JSONResponse
	checkError(json.Unmarshal(resp.Body, &token))

	return token.Token
}

func HeaderCreator(ctx context.Context, req *http.Request) error {
	req.Header.Add("Authorization", generateToken(ctx))
	return nil
}
