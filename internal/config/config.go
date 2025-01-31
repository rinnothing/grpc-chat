package config

import (
	"flag"
	"net"
)

const myUsername = "rinnothing"

const port = "8081"

func MustGetUsername() string {
	return myUsername
}

func MustGetPort() string {
	return port
}

func MustGetIPv4() net.IP {
	// copied from stack overflow
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP
			}
		}
	}
	return nil
}

// if true all incoming connections are accepted automatically
var acceptAll bool

func MustGetAcceptAll() bool {
	return acceptAll
}

func init() {
	flag.CommandLine.BoolVar(&acceptAll, "accept-all", false, "accept all incoming connections")
}
