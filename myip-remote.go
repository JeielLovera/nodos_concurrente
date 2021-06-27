package main

import (
	"fmt"
	"net"
	"strings"
)

func localAddress() string {
	ifaces, err := net.Interfaces()

	if err != nil {
		fmt.Printf("error")
	}
	fmt.Println(ifaces)
	for _, oiface := range ifaces {
		if strings.HasPrefix(oiface.Name, "ens33") {
			addrs, err := oiface.Addrs()

			if err != nil {
				fmt.Printf("error")
				continue
			}

			for _, dir := range addrs {
				switch d := dir.(type) {
				case *net.IPNet:
					if strings.HasPrefix(d.IP.String(), "192") {
						return d.IP.String()
					}
				}
			}
		}
	}
	return "127.0.0.1" // ah?
}

func main() {

}
