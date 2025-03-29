package test

import (
	"os"

	"github.com/AugustineAurelius/fuufu/internal/config"
	"github.com/AugustineAurelius/fuufu/pkg/logger"
	"github.com/spf13/cobra"
)

func checkError(err error) {
	if err != nil {
		logger.New().Error(err.Error())
		os.Exit(1)
	}
}

func CreateCMD(manager *config.Manager) *cobra.Command {

	createClients(manager)

	testCMD := &cobra.Command{
		Use: "test",
	}

	testCMD.AddCommand(createAuthCMD(manager), createTodoCMD(manager))
	return testCMD
}
