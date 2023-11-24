package main

import (
	"gontp/pkg/ntplib"
	"log"
	"net"
	"time"
)

func responseNTP(conn *net.UDPConn, addr *net.UDPAddr, requestData []byte) {
	recvTime := getPreciseCurrTime()

	request := ntplib.NTPMsg{}
	err := request.Decode(requestData)

	if err != nil {
		log.Print(err)
		return
	}

	_, ntpVn, ntpMode := request.NtpHead.GetHead()
	if ntpVn != 3 || ntpMode != 3 {
		log.Print("invalid ntp request")
		return
	}

	response := ntplib.NTPMsg{
		NtpHead:          ntplib.CreateNTPMsgHead(0, 3, 4),
		NtpStratum:       2,
		NtpOriTimestamp:  request.NtpOriTimestamp,
		NtpRecvTimestamp: recvTime,
	}

	sendTime := getPreciseCurrTime()
	response.NtpTransTimestamp = sendTime

	responseData, err := response.Encode()

	if err != nil {
		log.Print(err)
		return
	}

	_, err = conn.WriteToUDP(responseData, addr)
	if err != nil {
		log.Print(err)
		return
	}
}

// this is just a mock
func getPreciseCurrTime() ntplib.LongFixedPoint {
	timeNow := time.Now()
	return ntplib.CreateLongFixedPoint(timeNow)
}
