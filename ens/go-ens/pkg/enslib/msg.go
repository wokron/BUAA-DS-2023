package enslib

import (
	"bytes"
	"encoding/binary"
)

type ENSType byte

const (
	SUBSCRIBE ENSType = iota
	UNSUBSCRIBE
	PUBLISH
	UPDATE
)

const (
	MAX_TOPIC_LENGTH   int = 50
	MAX_MESSAGE_LENGTH int = 50
)

type ENSMsgData struct {
	Type    ENSType
	Topic   [MAX_TOPIC_LENGTH]byte
	Message [MAX_MESSAGE_LENGTH]byte
}

func (msgData *ENSMsgData) ToMsg() (msg *ENSMsg) {
	return &ENSMsg{
		Type:    msgData.Type,
		Topic:   string(msgData.Topic[:]),
		Message: string(msgData.Message[:]),
	}
}

type ENSMsg struct {
	Type    ENSType
	Topic   string
	Message string
}

func (msg *ENSMsg) ToMsgData() (msgData *ENSMsgData) {
	msgData = &ENSMsgData{Type: msg.Type}
	copy(msgData.Topic[:], []byte(msg.Topic))
	copy(msgData.Message[:], []byte(msg.Message))
	return msgData
}

func DecodeENSMsg(data []byte) (msg *ENSMsg, err error) {
	buf := bytes.Buffer{}

	_, err = buf.Write(data)
	if err != nil {
		return
	}

	msgData := &ENSMsgData{}

	err = binary.Read(&buf, binary.BigEndian, msgData)
	if err != nil {
		return
	}

	msg = msgData.ToMsg()

	return msg, nil
}

func (msg *ENSMsg) Encode() (data []byte, err error) {
	buf := bytes.Buffer{}
	msgData := msg.ToMsgData()

	err = binary.Write(&buf, binary.BigEndian, msgData)
	if err != nil {
		return
	}

	return buf.Bytes(), nil
}
