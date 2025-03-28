package test

import (
	"github.com/AugustineAurelius/fuufu/internal/config"
	"github.com/spf13/cobra"
)

func CreateCMD(manager *config.Manager) *cobra.Command {
	testCMD := &cobra.Command{
		Use: "test",
	}

	testCMD.AddCommand(createAuthCMD(manager))
	return testCMD
}
