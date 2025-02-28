package cmd

import (
	dnsdb "github.com/sohamjoshi25/dns-server/internal/dnsdb"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a DNS record",
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		rtype, _ := cmd.Flags().GetString("type")
		answer, _ := cmd.Flags().GetString("answer")

		if domain == "" || answer == "" {
			cmd.Println("Error: Both --domain and --answer flags are required")
			cmd.Usage()
			return
		}

		dnsdb.InsertRecord(domain, rtype, answer)
	},
}

func init() {
	addCmd.Flags().String("domain", "", "Domain name")
	addCmd.Flags().String("type", "A", "Record type")
	addCmd.Flags().String("answer", "", "Record answer")

	rootCmd.AddCommand(addCmd)
}
