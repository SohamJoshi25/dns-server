package cmd

import (
	dnsdb "github.com/sohamjoshi25/dns-server/internal/dnsdb"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all DNS records",
	Run: func(cmd *cobra.Command, args []string) {
		dnsdb.GetAllRecords()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
