package network

import (
	"fmt"
	"net/netip"
	"os/exec"
	"strings"
	"unsafe"

	"github.com/rs/zerolog/log"
	"github.com/zer0go/netguard-client/internal/config"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
	"golang.zx2c4.com/wireguard/windows/driver"
)

var (
	modWireGuard                 = windows.NewLazyDLL("wireguard.dll")
	procWireGuardCreate          = modWireGuard.NewProc("WireGuardCreateAdapter")
	procWireGuardClose           = modWireGuard.NewProc("WireGuardCloseAdapter")
	procWireGuardSetAdapterState = modWireGuard.NewProc("WireGuardSetAdapterState")
)

func (i *Interface) Create() error {
	interfaceName := config.Get().InterfaceName
	namePtr, _ := windows.UTF16PtrFromString(interfaceName)
	typePtr, _ := windows.UTF16PtrFromString("WireGuard")

	_, _, err := procWireGuardCreate.Call(
		uintptr(unsafe.Pointer(namePtr)),
		uintptr(unsafe.Pointer(typePtr)),
		0,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create WireGuard adapter")
	}

	return nil
}

func (i *Interface) LoadLink() error {
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
