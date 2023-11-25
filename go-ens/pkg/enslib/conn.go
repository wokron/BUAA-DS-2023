package enslib

import (
	"encoding/binary"
	"net"
)

func SendENSMsg(conn net.Conn, msg *ENSMsg) (err error) {
	sendData, err := msg.Encode()
	if err != nil {
		return
	}

	_, err = conn.Write(sendData)
	if err != nil {
		return
	}

	return nil
}

func RecvENSMsg(conn net.Conn) (msg *ENSMsg, err error) {
	dataSize := binary.Size(ENSMsgData{})

	recvData := make([]byte, dataSize)

	_, err = conn.Read(recvData)
	if err != nil {
		return
	}

	msg, err = DecodeENSMsg(recvData)
	if err != nil {
		return
	}

	return msg, nil
}
