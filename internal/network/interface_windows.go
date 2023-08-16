package network

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/netip"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows"
	"golang.zx2c4.com/wireguard/windows/driver"
)

func (i *Interface) Create() error {
	wgMutex.Lock()
	defer wgMutex.Unlock()

	err := i.LoadLink()
	if err != nil {
		log.Debug().Msg("creating windows tunnel")
		windowsGUID, err := windows.GenerateGUID()
		if err != nil {
			return err
		}
		i.Link, err = driver.CreateAdapter(i.Name, "WireGuard", &windowsGUID)
		if err != nil {
			return err
		}
	} else {
		log.Debug().Msg("re-using existing adapter")
	}

	log.Debug().Msg("created windows tunnel")

	return i.Link.(*driver.Adapter).SetAdapterState(driver.AdapterStateUp)
}

func (i *Interface) LoadLink() error {
	adapter, err := driver.OpenAdapter(i.Name)
	if err != nil {
		return err
	}
	i.Link = adapter

	return nil
}

func (i *Interface) ApplyMTU() error {
	return nil
}

func (i *Interface) ApplyAddress() error {
	var prefixAddresses []netip.Prefix
	for index := range i.Addresses {
		maskSize, _ := i.Addresses[index].Network.Mask.Size()
		log.Debug().
			Str("address", fmt.Sprintf("%s/%d to netguard interface", i.Addresses[index].IP.String(), maskSize)).
			Msg("appending address")

		addr, err := netip.ParsePrefix(fmt.Sprintf("%s/%d", i.Addresses[index].IP.String(), maskSize))
		if err == nil {
			prefixAddresses = append(prefixAddresses, addr)
		} else {
			log.Error().Err(err).Msg("failed to append ip to adapter")
		}
	}

	return i.Link.(*driver.Adapter).LUID().SetIPAddresses(prefixAddresses)
}

func (i *Interface) Close() error {
	err := i.Link.Close()
	if err != nil {
		log.Error().Err(err).Msg("error closing interface")
	}

	for index := range i.Addresses {
		if i.Addresses[index].Network.String() == "0.0.0.0/0" ||
			i.Addresses[index].Network.String() == "::/0" {
			continue
		}
		if i.Addresses[index].AddRoute {
			maskSize, _ := i.Addresses[index].Network.Mask.Size()
			log.Debug().
				Str("range", fmt.Sprintf("%s/%d from interface", i.Addresses[index].IP.String(), maskSize)).
				Msg("removing egress range")

			cmd := fmt.Sprintf("route delete %s", i.Addresses[index].IP.String())
			_, err := runCmd(cmd)
			if err != nil {
				log.Error().
					Str("address", i.Addresses[index].IP.String()).
					Msg("failed to remove egress range")
			}
		}
	}

	return nil
}

func runCmd(command string) (string, error) {
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)
	err := cmd.Wait()
	if err != nil {
		return "", err
	}

	out, err := cmd.CombinedOutput()

	return string(out), err
}
