package config

import "net"

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
