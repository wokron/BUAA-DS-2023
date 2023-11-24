package main

import (
	"goens/pkg/enslib/client"
	"log"
	"strconv"
)

func init() {
	flagInit()
}

func main() {
	serverAddr := serverName + ":" + strconv.Itoa(port)
	log.Printf("Connecting to ENS Server %s.\n", serverAddr)

	conn, err := client.CreateENSConn(serverAddr)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	for _, event := range events {
		log.Printf("Publish event on topic \"%s\", message: \"%s\"\n", event.topic, event.message)
		err = client.Publish(conn, event.topic, event.message)
		if err != nil {
			log.Printf("Error when publishing events, %s.\n", err)
		}
	}

	log.Println("All events has been published, exit.")
}
