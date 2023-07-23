package network

import (
	"github.com/rs/zerolog/log"
	"net"
)

type InterfaceAddress struct {
	IP       net.IP
	Network  net.IPNet
	AddRoute bool
}

func CreateFromCIDR(cidr string) []InterfaceAddress {
	var address []InterfaceAddress

	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Err(err)
		return address
	}
	address = append(address, InterfaceAddress{
		IP:      ip,
		Network: *ipNet,
	})

	return address
}
