package main

import (
	"gontp/pkg/ntplib"
	"net"
	"time"
)

func requestNTPServer(address string, timeoutSec int) (sendTime time.Time, recvTime time.Time, response ntplib.NTPMsg, err error) {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Duration(timeoutSec) * time.Second))

	return requestNTP(conn)
}

func requestNTP(conn net.Conn) (sendTime time.Time, recvTime time.Time, response ntplib.NTPMsg, err error) {
	request := ntplib.NTPMsg{
		NtpHead: ntplib.CreateNTPMsgHead(0, 3, 3),
		NtpPoll: 4,
	}

	responseData := make([]byte, 48)

	// >>> NTP start >>>
	sendTime = time.Now()
	request.NtpOriTimestamp = ntplib.CreateLongFixedPoint(sendTime)

	requestData, err := request.Encode()
	if err != nil {
		return
	}

	_, err = conn.Write(requestData)
	if err != nil {
		return
	}

	_, err = conn.Read(responseData)
	if err != nil {
		return
	}

	recvTime = time.Now()
	// <<< NTP end <<<

	response.Decode(responseData)
	return sendTime, recvTime, response, nil
}
