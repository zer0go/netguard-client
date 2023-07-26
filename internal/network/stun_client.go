package network

import (
	"github.com/pion/stun"
	"github.com/rs/zerolog/log"
	"net"
	"strings"
)

const (
	PublicNat = "public"
	BehindNat = "behind_nat"
)

var stunServerList = []string{
	"stun1.l.google.com:19302",
	"stun2.l.google.com:19302",
	"stun3.l.google.com:19302",
	"stun4.l.google.com:19302",
}

type StunClient struct {
	ServerList []string
}

func NewStunClientWithServers() *StunClient {
	return &StunClient{
		ServerList: stunServerList,
	}
}

func (s *StunClient) StunPort(port int) (publicIP net.IP, publicPort int, natType string, err error) {
	for _, stunServer := range s.ServerList {
		remoteAddress, err := net.ResolveUDPAddr("udp", stunServer)
		if err != nil {
			log.Error().Msgf("failed to resolve udp addr: %s %s", stunServer, err.Error())
			continue
		}
		listenAddress := &net.UDPAddr{
			IP:   net.ParseIP(""),
			Port: port,
		}

		publicIP, publicPort, natType, err = s.doStunTransaction(listenAddress, remoteAddress)
		if err != nil {
			log.Error().Msgf("stun transaction failed: %s %s", stunServer, err.Error())
			continue
		}
		if publicPort == 0 || publicIP == nil || publicIP.IsUnspecified() {
			continue
		}
		break
	}

	return
}

func (s *StunClient) doStunTransaction(listenAddress, remoteAddress *net.UDPAddr) (publicIP net.IP, publicPort int, natType string, err error) {
	conn, err := net.DialUDP("udp", listenAddress, remoteAddress)
	if err != nil {
		log.Error().Msgf("failed to dial: %s", err.Error())
		return
	}

	privateIP := net.ParseIP(extractIPFromAddress(conn.LocalAddr().String()))
	defer func() {
		if privateIP.Equal(publicIP) {
			natType = PublicNat
		} else {
			natType = BehindNat
		}
	}()
	defer func(conn *net.UDPConn) {
		_ = conn.Close()
	}(conn)

	c, err := stun.NewClient(conn)
	if err != nil {
		log.Error().Msgf("failed to create stun client: %s", err.Error())
		return
	}
	defer func(c *stun.Client) {
		_ = c.Close()
	}(c)

	if err = c.Do(stun.MustBuild(stun.TransactionID, stun.BindingRequest), func(res stun.Event) {
		if res.Error != nil {
			log.Error().Msgf("failed STUN transaction: %s", res.Error)
		}

		var xorAddr stun.XORMappedAddress
		if getErr := xorAddr.GetFrom(res.Message); getErr != nil {
			log.Error().Msgf("failed to get XOR-MAPPED-ADDRESS: %s", getErr)
		}

		publicIP = xorAddr.IP
		publicPort = xorAddr.Port
	}); err != nil {
		log.Error().Err(err).Msg("stun request error")
	}
	if err := c.Close(); err != nil {
		log.Error().Msgf("Failed to close connection: %s", err)
	}

	return
}

func extractIPFromAddress(address string) string {
	return strings.Split(address, ":")[0]
}
