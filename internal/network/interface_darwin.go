package network

import (
	"fmt"
	"net"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (i *Interface) Create() error {
	return nil
}

func (i *Interface) LoadLink() error {
	return nil
}

func (i *Interface) ApplyMTU() error {
	cmd := exec.Command("ifconfig", i.Name, "mtu", fmt.Sprintf("%d", i.MTU))
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to set MTU")
	}
	return nil
}

func (i *Interface) ApplyAddress() error {
	out, err := exec.Command("ifconfig", i.Name).CombinedOutput()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to get interface: %s", out))
	}

	// Parse and delete existing inet/inet6 addresses
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "inet ") || strings.HasPrefix(line, "inet6 ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				cmd := exec.Command("ifconfig", i.Name, strings.TrimSuffix(parts[0], ":"), parts[1], "delete")
				_ = cmd.Run()
			}
		}
	}

	// Add new addresses
	for _, addr := range i.Addresses {
		if addr.IP == nil || addr.Network.IP == nil {
			continue
		}

		if addr.IP.To4() != nil {
			// IPv4 utun (point-to-point)
			cmd := exec.Command(
				"ifconfig",
				i.Name,
				"inet",
				addr.IP.String(),         // local
				addr.Network.IP.String(), // destination (peer)
				"netmask",
				net.IP(addr.Network.Mask).String(),
			)

			if out, err := cmd.CombinedOutput(); err != nil {
				log.Error().
					Err(err).
					Str("output", string(out)).
					Msg("error adding IPv4 address")
			}
		} else {
			// IPv6 utun
			cidr := fmt.Sprintf("%s/%d", addr.IP.String(), maskToCIDR(addr.Network.Mask))
			cmd := exec.Command(
				"ifconfig",
				i.Name,
				"inet6",
				cidr,
			)

			if out, err := cmd.CombinedOutput(); err != nil {
				log.Error().
					Err(err).
					Str("output", string(out)).
					Msg("error adding IPv6 address")
			}
		}
	}

	err = i.applyRoutes()
	if err != nil {
		return errors.Wrap(err, "failed to apply routes")
	}

	return nil
}

func (i *Interface) applyRoutes() error {
	ipNet := i.Addresses[0].Network
	cidr := ipNetToCIDR(ipNet)
	peerIP := ipNet.IP.String()

	cmd := exec.Command(
		"route",
		"-n",
		"add",
		"-net",
		cidr,
		peerIP,
	)

	if out, err := cmd.CombinedOutput(); err != nil {
		log.Error().
			Err(err).
			Str("output", string(out)).
			Str("cidr", cidr).
			Str("peerIP", peerIP).
			Msg("failed to add route")
		return errors.Wrap(err, "failed to add route")
	}

	log.Debug().
		Str("cidr", cidr).
		Str("peerIP", peerIP).
		Msg("route added")

	return nil
}

func (i *Interface) Close() error {
	return nil
}

func ipNetToCIDR(ipNet net.IPNet) string {
	return fmt.Sprintf("%s/%d", ipNet.IP.String(), maskToCIDR(ipNet.Mask))
}

func maskToCIDR(mask net.IPMask) int {
	ones, _ := mask.Size()
	return ones
}
