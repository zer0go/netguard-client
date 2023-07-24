package network

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/vishvananda/netlink"
	"net"
	"os"
)

func (i *Interface) Create() error {
	err := i.LoadLink()
	if err != nil {
		return err
	}

	if err := i.deleteExistingLink(); err != nil {
		return err
	}

	if err := netlink.LinkAdd(i.Link.(netlink.Link)); err != nil && !os.IsExist(err) {
		return err
	}

	if err := netlink.LinkSetUp(i.Link.(netlink.Link)); err != nil {
		return err
	}

	return nil
}

func (i *Interface) LoadLink() error {
	link, err := i.getKernelLink()
	if err != nil {
		return err
	}
	i.Link = link

	return nil
}

func (i *Interface) ApplyMTU() error {
	link, err := i.getKernelLink()
	if err != nil {
		return err
	}
	if err := netlink.LinkSetMTU(link, i.MTU); err != nil {
		return err
	}

	return nil
}

func (i *Interface) ApplyAddress() error {
	link, err := netlink.LinkByName(i.Name)
	if err != nil {
		return errors.Errorf("failed to locate link %s", err)
	}

	routes, err := netlink.RouteList(link, 0)
	if err != nil {
		return err
	}
	currentAddrs, err := netlink.AddrList(link, 0)
	if err != nil {
		return err
	}

	for i := range routes {
		err = netlink.RouteDel(&routes[i])
		if err != nil {
			return errors.Errorf("failed to list routes %s", err)
		}
	}

	if len(currentAddrs) > 0 {
		for i := range currentAddrs {
			err = netlink.AddrDel(link, &currentAddrs[i])
			if err != nil {
				return errors.Errorf("failed to delete route %s", err)
			}
		}
	}

	for _, addr := range i.Addresses {
		if addr.IP != nil && addr.Network.IP != nil {
			log.Info().
				Str("address", addr.IP.String()).
				Str("network", addr.Network.String()).
				Msg("adding address")
			netLinkAddr := netlink.Addr{
				IPNet: &net.IPNet{
					IP:   addr.IP,
					Mask: addr.Network.Mask,
				},
			}
			if err := netlink.AddrAdd(link, &netLinkAddr); err != nil {
				log.Error().Err(err).Msg("error adding addrress")
			}
		}

	}

	return nil
}

func (i *Interface) Close() error {
	link, err := i.getKernelLink()
	if err != nil {
		return err
	}

	return link.Close()
}

func (i *Interface) getKernelLink() (*NetLink, error) {
	link := createKernelLink(i.Name)
	if link == nil {
		return nil, errors.New("failed to create kernel interface")
	}

	return link, nil
}

func (i *Interface) deleteExistingLink() error {
	existingLink, err := netlink.LinkByName(i.Name)
	if err != nil {
		switch err.(type) {
		case netlink.LinkNotFoundError:
			break
		default:
			return err
		}
	}

	if existingLink != nil {
		err = netlink.LinkDel(existingLink)
		if err != nil {
			return err
		}
	}

	return nil
}

func createKernelLink(name string) *NetLink {
	linkAttrs := netlink.NewLinkAttrs()
	linkAttrs.Name = name

	return &NetLink{
		attrs: &linkAttrs,
	}
}
