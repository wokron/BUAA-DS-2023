package main

import (
	"flag"
	"fmt"
	"strconv"
	"syscall"
	"time"
)

func main() {
	flag.Parse()

	sendTime, recvTime, response, err := requestNTPServer(serverName+":"+strconv.Itoa(port), timeout)
	if err != nil {
		fmt.Println("fail to request ntp,", err)
		return
	}

	rrt := response.CalcRRT(sendTime, recvTime)
	offset := response.CalcOffset(sendTime, recvTime)

	if isShowResult {
		fmt.Println("client send:", sendTime)
		fmt.Println("server recv:", response.NtpRecvTimestamp.ToTime())
		fmt.Println("server send:", response.NtpTransTimestamp.ToTime())
		fmt.Println("client recv:", recvTime)
		fmt.Println()
		fmt.Println("rrt:", rrt, "(µs)")
		fmt.Println("offset:", offset, "(µs)")
	}

	if isSetSysTime {
		err := setSysTime(time.Now().Add(time.Duration(offset) * time.Microsecond))
		if err != nil {
			fmt.Println("fail to set system time,", err)
		}
	}
}

func setSysTime(ntpTime time.Time) error {
	err := syscall.Settimeofday(&syscall.Timeval{Sec: ntpTime.Unix()})
	if err != nil {
		return err
	}
	return nil
}
