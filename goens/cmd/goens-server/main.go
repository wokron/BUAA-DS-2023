package main

import (
	"goens/pkg/enslib"
	"goens/pkg/enslib/server"
	"io"
	"log"
	"net"
)

func handleConnection(conn net.Conn, table *server.SubscribeTable) {
	defer conn.Close()
	log.Println("tcp connect success")

	for {
		msg, err := enslib.RecvENSMsg(conn)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
		}

		switch msg.Type {
		case enslib.SUBSCRIBE:
			table.Subscribe(msg.Topic, conn)
		case enslib.UNSUBSCRIBE:
			table.Unsubscribe(msg.Topic, conn)
		case enslib.PUBLISH:
			conns := table.GetSubscribers(msg.Topic)
			server.PublishEventToAll(conns, msg.Topic, msg.Message, func(err error) { log.Println(err) })
		default:
			log.Printf("msg type unknown")
		}
	}
	log.Printf("conn %s has been closed\n", conn.RemoteAddr().String())
}

func main() {
	table := server.NewSubscribeTable()

	listen, err := net.Listen("tcp", "0.0.0.0:4567")
	if err != nil {
		log.Panic(err)
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConnection(conn, table)
	}
}
