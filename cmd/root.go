package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:     "dns",
	Short:   "A CLI App which allows to start and have own DNS server with resolver and having ability to add own DNS records.",
	Long:    "A CLI App which allows to start and have own DNS server with resolver and having ability to add own DNS records.",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to DNS CLI")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Optional: Add short version flag (-v)
	rootCmd.Flags().BoolP("version", "v", false, "Print the version number")
}
