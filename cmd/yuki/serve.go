package yuki

import (
	"github.com/AugustineAurelius/fuufu/internal/config"
	"github.com/AugustineAurelius/fuufu/yuki"

	"github.com/spf13/cobra"
)

func createServeCMD(manager *config.Manager) *cobra.Command {
	serveCMD := &cobra.Command{
		Use: "serve",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := manager.LoadYuki()
			yuki.Run(cfg.Host, cfg.Port)
		},
	}

	return serveCMD
}
