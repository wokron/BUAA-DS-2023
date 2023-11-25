package main

import (
	"flag"
	"log"
	"net"
	"strconv"
)

func main() {
	flag.Parse()

	udpAddr, err := net.ResolveUDPAddr("udp", host+":"+strconv.Itoa(port))
	if err != nil {
		log.Panic(err)
	}

	log.Print("udp address: ", udpAddr)

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Panic(err)
	}

	defer conn.Close()

	buf := make([]byte, 48)

	log.Print("start listening ntp requests...")
	for {
		_, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("recv request from %s\n", addr)
		go responseNTP(conn, addr, buf)
	}
}
