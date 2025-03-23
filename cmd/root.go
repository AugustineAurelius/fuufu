/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/AugustineAurelius/fuufu/internal/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func Execute() {
	printWelcome()

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

func printWelcome() {
	f := []string{
		"███████╗",
		"██╔════╝",
		"█████╗  ",
		"██╔══╝  ",
		"██║     ",
	}

	u := []string{
		"██╗  ██╗",
		"██║  ██║",
		"██║  ██║",
		"╚█████╔╝",
		" ╚════╝ ",
	}

	pattern := []struct {
		char  string
		color *color.Color
	}{
		{"F", color.New(color.FgHiWhite)},
		{"U", color.New(color.FgHiWhite)},
		{"U", color.New(color.FgHiWhite)},
		{"F", color.New(color.FgHiWhite)},
		{"U", color.New(color.FgHiWhite)},
	}

	for i := range 5 {
		for _, p := range pattern {
			switch p.char {
			case "F":
				p.color.Print(f[i])
			case "U":
				p.color.Print(u[i])
			}
		}
		fmt.Println()
	}

	// Дополнительная декоративная линия
	color.HiWhite("✻" + strings.Repeat("─", 36) + "✻")
	color.HiWhite("  Welcome to our family application!")
}
