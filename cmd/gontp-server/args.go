package main

import "flag"

var (
	host string
	port int
)

func init() {
	flag.StringVar(&host, "host", "0.0.0.0", "the host name of ntp service")
	flag.IntVar(&port, "port", 123, "the port of ntp service")
}
