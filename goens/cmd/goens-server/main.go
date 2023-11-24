package main

import (
	"goens/pkg/enslib"
	"goens/pkg/enslib/server"
	"io"
	"log"
	"net"
	"strconv"
)

func init() {
	flagInit()
}

func handleConnection(conn net.Conn, table *server.SubscribeTable) {
	defer conn.Close()
	log.Printf("Accept tcp connection from %s.\n", conn.RemoteAddr().String())

	connAddr := conn.RemoteAddr().String()

	for {
		msg, err := enslib.RecvENSMsg(conn)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("%s: Error when receving data, %s.\n", connAddr, err)
		}

		switch msg.Type {
		case enslib.SUBSCRIBE:
			log.Printf("%s: Subscribe topic %s.\n", connAddr, msg.Topic)
			table.Subscribe(msg.Topic, conn)

		case enslib.UNSUBSCRIBE:
			log.Printf("%s: Unsubscribe topic %s.\n", connAddr, msg.Topic)
			table.Unsubscribe(msg.Topic, conn)

		case enslib.PUBLISH:
			log.Printf("%s: Publish event on topic %s, message: \"%s\".\n", connAddr, msg.Topic, msg.Message)
			conns := table.GetSubscribers(msg.Topic)
			server.PublishEventToAll(conns, msg.Topic, msg.Message, func(err error) { log.Println(err) })

		default:
			log.Printf("%s: Unknown Message Type.\n", connAddr)
		}
	}
	log.Printf("%s: Connection close.\n", connAddr)
}

func main() {
	table := server.NewSubscribeTable()

	addr := host + ":" + strconv.Itoa(port)
	log.Printf("Event Notification Service running on %s.\n", addr)

	listen, err := net.Listen("tcp", addr)
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
