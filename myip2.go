package main

import (
	"fmt"
	"net"
	"strings"
)

func zzmain() {
	fmt.Println(myIp())
}
func myIp() string { // mandrakeando ando
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		fmt.Println("entre1")
		fmt.Println(iface.Name)
		if strings.HasPrefix(iface.Name, "192") {
			fmt.Println("entre2")
			addrs, _ := iface.Addrs()
			for _, addr := range addrs {
				switch v := addr.(type) {
				case *net.IPNet:
					fmt.Println(v.IP.String())
					return v.IP.String()
				case *net.IPAddr:
					fmt.Println(v.IP.String())
					return v.IP.String()
				}
			}
		}
	}
	return "a"
}
