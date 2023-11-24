package main

import (
	"fmt"
	"goens/pkg/enslib"
	"goens/pkg/enslib/client"
	"io"
	"log"
	"net"
	"time"
)

func listenEvents(conn net.Conn, done chan struct{}) {
	for {
		select {
		case <-done:
			return
		default:
			conn.SetReadDeadline(time.Now().Add(1 * time.Microsecond))

			msg, err := enslib.RecvENSMsg(conn)
			if err == io.EOF {
				break
			} else if err != nil && err.(*net.OpError).Timeout() {
				continue
			} else if err != nil {
				log.Println(err)
				continue
			}

			fmt.Printf("%s: %s\n", msg.Topic, msg.Message)
		}
	}
}

func nextToken() string {
	var token string
	fmt.Scanf("%s", &token)
	return token
}

func main() {
	conn, err := client.CreateENSConn("localhost:4567")
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	listenDone := make(chan struct{})
	go listenEvents(conn, listenDone)

loop:
	for {
		cmd := nextToken()
		switch cmd {
		case "quit":
			listenDone <- struct{}{}
			break loop

		case "pub":
			fallthrough
		case "public":
			topic := nextToken()
			message := nextToken()
			client.Publish(conn, topic, message)

		case "sub":
			fallthrough
		case "subscribe":
			topic := nextToken()
			client.Subscribe(conn, topic)

		case "unsub":
			fallthrough
		case "unsubscribe":
			topic := nextToken()
			client.Unsubscribe(conn, topic)

		default:
			fmt.Println("cmd unknown:", cmd)
		}
	}
}
