package network

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/zer0go/netguard-client/internal/config"
	"sync"
)

var wgMutex = sync.Mutex{}

type Interface struct {
	Name      string
	Link      NetLinkInterface
	Addresses []InterfaceAddress
	MTU       int
}

type NetLinkInterface interface {
	Close() error
}

func NewInterfaceFromConfig(c *config.App) *Interface {
	return &Interface{
		Name:      c.InterfaceName,
		MTU:       c.MTU,
		Addresses: getAddressesFromCIDR(c.NetworkCIDR),
	}
}

func (i *Interface) Configure() error {
	wgMutex.Lock()
	defer wgMutex.Unlock()

	log.Debug().Msg("adding addresses to interface")

	if err := i.ApplyMTU(); err != nil {
		return errors.Errorf("configure set MTU %s", err)
	}
	if err := i.ApplyAddress(); err != nil {
		return err
	}

	return nil
}

func getAddressesFromCIDR(cidr string) []InterfaceAddress {
	var addresses []InterfaceAddress
	if cidr == "" {
		return addresses
	}

	interfaceAddress, err := CreateInterfaceAddressFromCIDR(cidr)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse network cidr")
		return addresses
	}

	return append(addresses, *interfaceAddress)
}
