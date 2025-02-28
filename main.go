package main

import (
	"github.com/sohamjoshi25/dns-server/cmd"
	_ "github.com/sohamjoshi25/dns-server/internal/dnsdb"
)

func main() {
	cmd.Execute()
}
