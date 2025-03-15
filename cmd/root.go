/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/AugustineAurelius/fuufu/internal/config"
	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := createRootCMD()

	manager := config.NewManager(rootCmd)

	rootCmd.AddCommand(createServeCMD(manager), createMigrateCMD(manager))

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func createRootCMD() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "fuufu",
		Short: "Our family app",
	}

	rootCmd.PersistentFlags().StringP("config", "c", "", "Path to a configuration file")
	return rootCmd
}
