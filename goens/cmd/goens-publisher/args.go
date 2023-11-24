package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type event struct {
	topic   string
	message string
}

type eventList []event

func (events *eventList) String() string {
	strEvents := []string{}
	for _, event := range *events {
		strEvents = append(strEvents, event.topic+":"+event.message)
	}
	return fmt.Sprintf("%v", strEvents)
}

func (events *eventList) Set(arg string) error {
	split := strings.SplitN(arg, ":", 2)

	if len(split) != 2 {
		return errors.New("false format, should be <topic>:<message>")
	}
	*events = append(*events, struct {
		topic   string
		message string
	}{topic: split[0], message: split[1]})

	return nil
}

var (
	serverName string
	port       int
	events     eventList
)

func flagInit() {
	flag.StringVar(&serverName, "server", "localhost", "address of ENS server")
	flag.IntVar(&port, "port", 4567, "the port of event notification sercice")
	flag.Var(&events, "event", "event you want to publish, the format should be <topic>:<message>")
	flag.Parse()
}
