package crosslab

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "crosslab",
	Short: "Crosslab is a CLI tool",
	Long: `A CLI application built with love using Cobra.
This is the root command for the application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
