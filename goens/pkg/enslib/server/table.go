package server

import (
	"log"
	"net"
	"sync"
)

type SubscribeTable struct {
	mutex sync.RWMutex
	subs  map[string][]net.Conn
}

func NewSubscribeTable() *SubscribeTable {
	table := &SubscribeTable{
		subs: make(map[string][]net.Conn),
	}
	return table
}

func (table *SubscribeTable) Subscribe(topic string, conn net.Conn) {
	table.mutex.Lock()
	defer table.mutex.Unlock()

	table.subs[topic] = append(table.subs[topic], conn)
}

func (table *SubscribeTable) Unsubscribe(topic string, conn net.Conn) {
	table.mutex.Lock()
	defer table.mutex.Unlock()

	findIdx := func(conns []net.Conn, targetConn net.Conn) int {
		for i, conn := range conns {
			if conn == targetConn {
				return i
			}
		}
		return -1
	}

	idx := findIdx(table.subs[topic], conn)
	if idx < 0 {
		return
	}

	table.subs[topic] = append(table.subs[topic][:idx], table.subs[topic][idx+1:]...)
}

func (table *SubscribeTable) GetSubscribers(topic string) (subscribers []net.Conn) {
	table.mutex.Lock()
	defer table.mutex.Unlock()

	subscribers = make([]net.Conn, len(table.subs[topic]))
	copy(subscribers, table.subs[topic])
	log.Print(len(subscribers))
	return
}
