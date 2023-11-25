package client

import (
	"goens/pkg/enslib"
	"net"
)

func CreateENSConn(address string) (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", address)
	if err != nil {
		return
	}
	return conn, nil
}

func Publish(conn net.Conn, topic, message string) (err error) {
	msg := &enslib.ENSMsg{
		Type:    enslib.PUBLISH,
		Topic:   topic,
		Message: message,
	}

	err = enslib.SendENSMsg(conn, msg)
	if err != nil {
		return
	}

	return nil
}

func Subscribe(conn net.Conn, topic string) (err error) {
	msg := &enslib.ENSMsg{
		Type:  enslib.SUBSCRIBE,
		Topic: topic,
	}

	err = enslib.SendENSMsg(conn, msg)
	if err != nil {
		return
	}

	return nil
}

func Unsubscribe(conn net.Conn, topic string) (err error) {
	msg := &enslib.ENSMsg{
		Type:  enslib.UNSUBSCRIBE,
		Topic: topic,
	}

	err = enslib.SendENSMsg(conn, msg)
	if err != nil {
		return
	}

	return nil
}
