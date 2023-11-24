package main

import "flag"

var (
	serverName   string
	port         int
	timeout      int
	isShowResult bool
	isSetSysTime bool
)

func init() {
	flag.StringVar(&serverName, "server", "ntp.ntsc.ac.cn", "the server of ntp")
	flag.IntVar(&port, "port", 123, "the port of ntp service")
	flag.IntVar(&timeout, "timeout", 10, "request timeout sec")
	flag.BoolVar(&isShowResult, "result", true, "whether to show the rrt and offset")
	flag.BoolVar(&isSetSysTime, "settime", false, "whether to modify the system time")
}
