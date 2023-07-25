package handler

import (
	"encoding/base64"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/zer0go/netguard-client/internal/config"
	"github.com/zer0go/netguard-client/internal/network"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"io"
	"net"
	"net/http"
	"time"
)

type Peer struct {
	AllowedIP string `json:"allowed_ip"`
	Endpoint  string `json:"endpoint"`
	PublicKey string `json:"public_key"`
}

type WgConfig struct {
	NetworkCIDR   string `json:"network_cidr"`
	ListeningPort int    `json:"listening_port"`
	PrivateKey    string `json:"private_key"`
	Peers         []Peer `json:"peers"`
}

type JoinHandler struct {
}

func NewJoinHandler() *JoinHandler {
	return &JoinHandler{}
}

func (h *JoinHandler) Handle(cmd *cobra.Command, _ []string) error {
	log.Info().Msg("joining to server...")

	token, _ := cmd.Flags().GetString("token")
	joinConfigJSON, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return err
	}
	var joinConfig map[string]string
	err = json.Unmarshal(joinConfigJSON, &joinConfig)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodGet, joinConfig["url"], nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+joinConfig["token"])

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	response, err := client.Do(req)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var wgConfig WgConfig
	err = json.Unmarshal(responseBody, &wgConfig)
	if err != nil {
		return err
	}

	log.Debug().
		Str("response_body", string(responseBody)).
		Msg("join api response received")

	appConfig := config.Get()
	appConfig.NetworkCIDR = wgConfig.NetworkCIDR
	err = config.Update(*appConfig)
	if err != nil {
		return err
	}

	appConfig = config.Get()
	log.Debug().
		Interface("AppConfig", appConfig).
		Msg("")

	networkInterface := network.NewInterfaceFromConfig(appConfig)
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
	privateKey, _ := wgtypes.ParseKey(wgConfig.PrivateKey)

	keepaliveInterval := time.Second * 20
	var peers []wgtypes.PeerConfig
	for _, peer := range wgConfig.Peers {
		var allowedIPs []net.IPNet
		ipNet := net.IPNet{
			IP:   net.ParseIP(peer.AllowedIP),
			Mask: net.IPv4Mask(255, 255, 255, 255),
		}
		allowedIPs = append(allowedIPs, ipNet)
		peerEndpoint, err := net.ResolveUDPAddr("udp", peer.Endpoint)
		if err != nil {
			return err
		}
		peerPublicKey, err := wgtypes.ParseKey(peer.PublicKey)
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
	}

	return wg.ConfigureDevice(config.Get().InterfaceName, wgtypes.Config{
		PrivateKey:   &privateKey,
		FirewallMark: &firewallMark,
		ListenPort:   &wgConfig.ListeningPort,
		Peers:        peers,
		ReplacePeers: true,
	})
}
