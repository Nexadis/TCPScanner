package main

import (
	"flag"
	"log"

	"github.com/Nexadis/TCPTools/internal/blocker"
)

var blocklist, endpoint string

func main() {
	flag.StringVar(&blocklist, "bl", "blocklist.txt", "Blocklist with list of blocked sites")
	flag.StringVar(&endpoint, "l", "localhost:8080", "Address for listening http requests")
	flag.Parse()
	firewall, err := blocker.New(blocklist, endpoint)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(firewall.Run())
}
