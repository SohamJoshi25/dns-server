package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	_ "github.com/lib/pq"
	dnslookup "github.com/sohamjoshi25/go-dns-server/dnslookup"
	dnsproxy "github.com/sohamjoshi25/go-dns-server/dnsproxy"
)

func main() {

	cache := expirable.NewLRU[dnslookup.DNSQuestion, []dnslookup.DNSAnswer](16, nil, time.Minute*5)

	fmt.Printf("DNS Server Running on port 53")

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
}
