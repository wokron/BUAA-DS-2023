package main

import "flag"

var (
	host string
	port int
)

func flagInit() {
	flag.StringVar(&host, "host", "0.0.0.0", "the host address of the server")
	flag.IntVar(&port, "port", 4567, "the port of event notification sercice")
	flag.Parse()
}
