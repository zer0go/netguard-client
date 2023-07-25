package network

import (
	"net"
)

type InterfaceAddress struct {
	IP       net.IP
	Network  net.IPNet
	AddRoute bool
}

func CreateInterfaceAddressFromCIDR(cidr string) (*InterfaceAddress, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	return &InterfaceAddress{
		IP:      ip,
		Network: *ipNet,
	}, nil
}
