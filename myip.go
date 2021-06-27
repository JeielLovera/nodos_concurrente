package main

import (
	"fmt"
	"net"
	"os"
)

func zzxmain() {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			fmt.Println(ipv4)
		}
	}
}
