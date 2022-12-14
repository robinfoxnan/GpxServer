package main

import (
	"fmt"
	"github.com/go-ping/ping"
)

func main() {
	pinger, err := ping.NewPinger("10.128.106.99")
	pinger.SetPrivileged(true)
	if err != nil {
		panic(err)
	}
	pinger.Count = 3
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		panic(err)
	}
	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	fmt.Println(stats)
}
