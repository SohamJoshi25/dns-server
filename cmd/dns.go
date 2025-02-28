package cmd

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	dnslookup "github.com/sohamjoshi25/dns-server/internal/dnslookup"
	dnsproxy "github.com/sohamjoshi25/dns-server/internal/dnsproxy"
	"github.com/spf13/cobra"
)

var dnsStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the DNS server",
	Run: func(cmd *cobra.Command, args []string) {
		cache := expirable.NewLRU[dnslookup.DNSQuestion, []dnslookup.DNSAnswer](16, nil, time.Minute*5)
		fmt.Println("DNS Server Running on port 53")

		addr := net.UDPAddr{
			Port: 53,
			IP:   net.ParseIP("0.0.0.0"),
		}

		conn, err := net.ListenUDP("udp", &addr)
		if err != nil {
			fmt.Println("Could not start server:", err)
			os.Exit(1)
		}
		defer conn.Close()

		for {
			dnsproxy.HandleDNSRequest(conn, cache)
		}
	},
}

func init() {
	rootCmd.AddCommand(dnsStartCmd)
}
