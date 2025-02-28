package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mycli",
	Short: "MyCLI is a simple CLI app",
	Long:  `MyCLI is a CLI built with Cobra in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to MyCLI")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
	}
}
