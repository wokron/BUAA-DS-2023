package main

import (
	"goens/pkg/enslib"
	"goens/pkg/enslib/client"
	"io"
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

	for _, topic := range topics {
		log.Printf("Subscrib topic \"%s\"\n", topic)
		err = client.Subscribe(conn, topic)
		if err != nil {
			log.Printf("Error when subscribing topics, %s.\n", err)
		}
	}

	for {
		msg, err := enslib.RecvENSMsg(conn)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("Received Event, topic: \"%s\", message: \"%s\"\n", msg.Topic, msg.Message)
	}
}
