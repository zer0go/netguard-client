package network

import "github.com/vishvananda/netlink"

const (
	NetlinkType = "wireguard"
)

type NetLink struct {
	attrs *netlink.LinkAttrs
}

func (l *NetLink) Close() error {
	return netlink.LinkDel(l)
}

func (l *NetLink) Attrs() *netlink.LinkAttrs {
	return l.attrs
}

func (l *NetLink) Type() string {
	return NetlinkType
}
