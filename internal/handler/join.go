package handler

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/zer0go/netguard-client/internal/config"
	"github.com/zer0go/netguard-client/internal/network"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"net"
	"time"
)

type JoinHandler struct {
}

func NewJoinHandler() *JoinHandler {
	return &JoinHandler{}
}

func (h *JoinHandler) Handle(cmd *cobra.Command, _ []string) error {
	log.Info().Msg("joining to server...")

	token, _ := cmd.Flags().GetString("token")

	log.Print(token)

	peerAllowedIP, _ := cmd.Flags().GetString("peer_allowed_ip")
	peerEndpointAddress, _ := cmd.Flags().GetString("peer_endpoint")
	peerPublicKeyBase64, _ := cmd.Flags().GetString("peer_public_key")

	networkInterface := network.NewInterfaceFromConfig(config.Get())
	if err := networkInterface.Configure(); err != nil {
		return err
	}

	wg, err := wgctrl.New()
	if err != nil {
		return errors.Errorf("wgctrl %s", err)
	}
	defer func(wg *wgctrl.Client) {
		_ = wg.Close()
	}(wg)

	firewallMark := 0
	privateKey, _ := wgtypes.ParseKey(privateKeyBase64)
	log.Debug().
		Str("privateKey", privateKey.String()).
		Str("publicKey", privateKey.PublicKey().String()).
		Int("listeningPort", listeningPort).
		Str("peerAllowedIP", peerAllowedIP).
		Str("peerEndpointAddress", peerEndpointAddress).
		Str("peerPublicKeyBase64", peerPublicKeyBase64).
		Msg("wireguard device configured")

	var peers []wgtypes.PeerConfig
	keepaliveInterval := time.Second * 20
	var allowedIPs []net.IPNet
	ipNet := net.IPNet{
		IP:   net.ParseIP(peerAllowedIP),
		Mask: net.IPv4Mask(255, 255, 255, 255),
	}
	allowedIPs = append(allowedIPs, ipNet)
	peerEndpoint, err := net.ResolveUDPAddr("udp", peerEndpointAddress)
	if err != nil {
		return err
	}
	peerPublicKey, err := wgtypes.ParseKey(peerPublicKeyBase64)
	if err != nil {
		return err
	}

	peers = append(peers, wgtypes.PeerConfig{
		AllowedIPs:                  allowedIPs,
		Endpoint:                    peerEndpoint,
		PublicKey:                   peerPublicKey,
		PersistentKeepaliveInterval: &keepaliveInterval,
		//Remove:                      false,
		//UpdateOnly:                  false,
		//ReplaceAllowedIPs:           false,
		//PresharedKey:                nil,
	})

	return wg.ConfigureDevice(config.Get().InterfaceName, wgtypes.Config{
		PrivateKey:   &privateKey,
		FirewallMark: &firewallMark,
		ListenPort:   &listeningPort,
		Peers:        peers,
		ReplacePeers: true,
	})
}
