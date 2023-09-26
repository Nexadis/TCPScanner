package config

import "flag"

type Config struct {
	BlackList string
	WhiteList string
	Address   string
}

func New() *Config {
	c := &Config{}
	flag.StringVar(&c.BlackList, "bl", "", "Blacklist with list of blocked sites")
	flag.StringVar(&c.WhiteList, "wl", "", "Whitelist with list of allowed sites")
	flag.StringVar(&c.Address, "l", "localhost:8080", "Address for listening http requests")
	c.Parse()
	return &Config{}
}

func (c Config) Parse() {
	flag.Parse()
}
