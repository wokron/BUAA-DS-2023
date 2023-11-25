package main

import (
	"flag"
	"fmt"
)

type topicSlice []string

func (topics *topicSlice) String() string {
	return fmt.Sprintf("%v", []string(*topics))
}

func (topics *topicSlice) Set(topic string) error {
	*topics = append(*topics, topic)
	return nil
}

var (
	serverName string
	port       int
	topics     topicSlice
)

func flagInit() {
	flag.StringVar(&serverName, "server", "localhost", "Address of ENS server")
	flag.IntVar(&port, "port", 4567, "The port of event notification sercice")
	flag.Var(&topics, "topic", "Topic of event you want to subscribe")
	flag.Parse()
}
