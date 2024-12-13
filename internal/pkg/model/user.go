package model

import "net"

type User struct {
	ID       int
	Username string
	IPv4     net.IP
}
