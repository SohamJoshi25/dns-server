package cmd

import (
	dnsdb "github.com/sohamjoshi25/dns-server/internal/dnsdb"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a DNS record by ID",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetInt("id")

		if id <= 0 {
			cmd.Println("Error: Both --id flags are required ans cannot be less than or equal to 0")
			cmd.Usage()
			return
		}

		dnsdb.DeleteRecordByID(id)
	},
}

func init() {
	deleteCmd.Flags().Int("id", 0, "Record ID to delete")
	rootCmd.AddCommand(deleteCmd)
}
