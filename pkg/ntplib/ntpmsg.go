package ntplib

import (
	"bytes"
	"encoding/binary"
	"time"
)

type NTPMsgHead uint8

func CreateNTPMsgHead(ntpLi, ntpVn, ntpMode uint8) NTPMsgHead {
	return (NTPMsgHead)((ntpLi << 6) | ((ntpVn & 0b111) << 3) | (ntpMode & 0b111))
}

func (head *NTPMsgHead) GetHead() (ntpLi uint8, ntpVn uint8, ntpMode uint8) {
	data := (uint8)(*head)
	ntpLi = (data & 0b11000000) >> 6
	ntpVn = (data & 0b00111000) >> 3
	ntpMode = (data & 0b00000111)
	return
}

type NTPMsg struct {
	NtpHead           NTPMsgHead
	NtpStratum        uint8
	NtpPoll           uint8
	NtpPrecission     uint8
	NtpRootDelay      ShortFixedPoint
	NtpRootDispersion ShortFixedPoint
	NtpRefIdentifier  uint32
	NtpRefTimestamp   LongFixedPoint
	NtpOriTimestamp   LongFixedPoint
	NtpRecvTimestamp  LongFixedPoint
	NtpTransTimestamp LongFixedPoint
}

func (ntpMsg *NTPMsg) Decode(data []byte) (err error) {
	buf := bytes.Buffer{}

	_, err = buf.Write(data)
	if err != nil {
		return
	}

	err = binary.Read(&buf, binary.BigEndian, ntpMsg)
	if err != nil {
		return
	}

	return nil
}

func (ntpMsg *NTPMsg) Encode() (data []byte, err error) {
	buf := bytes.Buffer{}

	err = binary.Write(&buf, binary.BigEndian, ntpMsg)
	if err != nil {
		return
	}

	return buf.Bytes(), nil
}

func (ntpMsg *NTPMsg) CalcRRT(sendTime, recvTime time.Time) int64 {
	t1 := sendTime.UnixMicro()
	t2 := ntpMsg.NtpRecvTimestamp.ToUnixMicro()
	t3 := ntpMsg.NtpTransTimestamp.ToUnixMicro()
	t4 := recvTime.UnixMicro()

	return (t4 - t1) - (t3 - t2)
}

func (ntpMsg *NTPMsg) CalcOffset(sendTime, recvTime time.Time) int64 {
	t1 := sendTime.UnixMicro()
	t2 := ntpMsg.NtpRecvTimestamp.ToUnixMicro()
	t3 := ntpMsg.NtpTransTimestamp.ToUnixMicro()
	t4 := recvTime.UnixMicro()

	return ((t2 - t1) + (t3 - t4)) / 2
}
