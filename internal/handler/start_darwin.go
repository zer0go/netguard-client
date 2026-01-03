package handler

import (
	"net"

	"github.com/rs/zerolog/log"
	"github.com/zer0go/netguard-client/internal/config"
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/ipc"
	"golang.zx2c4.com/wireguard/tun"
)

func startNetworkAdapter() (err error, closeFunc func()) {
	log.Debug().
		Str("interface", config.Get().InterfaceName).
		Int("mtu", config.Get().MTU).
		Msg("creating TUN device...")

	tunDev, err := tun.CreateTUN(config.Get().InterfaceName, config.Get().MTU)
	if err != nil {
		return err, nil
	}

	ifName, err := tunDev.Name()
	if err != nil {
		return err, nil
	}

	log.Debug().
		Str("tunInterface", ifName).
		Msg("created TUN device")

	logger := device.NewLogger(device.LogLevelSilent, "")
	wgDev := device.NewDevice(
		tunDev,
		conn.NewDefaultBind(),
		logger,
	)

	log.Debug().
		Msg("bringing up WireGuard device...")

	if err := wgDev.Up(); err != nil {
		return err, nil
	}

	log.Debug().
		Msg("WireGuard device is up")

	log.Debug().
		Msg("starting UAPI listener...")

	uapiListener, err := startUAPI(ifName, wgDev)
	if err != nil {
		return err, nil
	}

	log.Debug().
		Msg("UAPI listener started")

	return nil, func() {
		_ = uapiListener.Close()
		wgDev.Close()
		_ = tunDev.Close()
	}
}

func startUAPI(ifName string, wgDev *device.Device) (net.Listener, error) {
	uapiSock, err := ipc.UAPIOpen(ifName)
	if err != nil {
		return nil, err
	}

	listener, err := ipc.UAPIListen(ifName, uapiSock)
	if err != nil {
		_ = uapiSock.Close()
		return nil, err
	}

	go func() {
		for {
			c, err := listener.Accept()
			if err != nil {
				return
			}
			go wgDev.IpcHandle(c)
		}
	}()

	return listener, nil
}
