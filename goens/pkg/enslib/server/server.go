package server

import (
	"goens/pkg/enslib"
	"net"
)

func PublishENSMsgToAll(conns []net.Conn, msg *enslib.ENSMsg, errHandler func(error)) {
	for _, conn := range conns {
		err := enslib.SendENSMsg(conn, msg)
		if err != nil {
			errHandler(err)
		}
	}
}

func PublishEventToAll(conns []net.Conn, topic, message string, errHandler func(error)) {
	msg := &enslib.ENSMsg{
		Type:    enslib.UPDATE,
		Topic:   topic,
		Message: message,
	}

	PublishENSMsgToAll(conns, msg, errHandler)
}
