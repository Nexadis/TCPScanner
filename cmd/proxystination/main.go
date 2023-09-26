package main

import (
	"log"

	"github.com/Nexadis/TCPTools/internal/blocker"
	"github.com/Nexadis/TCPTools/internal/blocker/config"
)

func main() {
	c := config.New()
	firewall, err := blocker.New(c)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(firewall.Run())
}
