package yuki

import (
	"github.com/AugustineAurelius/fuufu/internal/config"
	"github.com/spf13/cobra"
)

func CreateCMD(manager *config.Manager) *cobra.Command {
	yukiCMD := &cobra.Command{
		Use: "yuki",
	}

	yukiCMD.AddCommand(createServeCMD(manager), registerTestCMD(manager))
	return yukiCMD
}
