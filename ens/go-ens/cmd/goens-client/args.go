package main

import "flag"

var (
	serverName string
	port       int
)

func flagInit() {
	flag.StringVar(&serverName, "server", "localhost", "Address of ENS server")
	flag.IntVar(&port, "port", 4567, "The port of event notification sercice")
	flag.Parse()
}
