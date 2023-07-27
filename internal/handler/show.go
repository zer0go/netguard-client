package handler

import (
	"bytes"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type ShowHandler struct {
}

func NewShowHandler() *ShowHandler {
	return &ShowHandler{}
}

func (h *ShowHandler) Handle(_ *cobra.Command, args []string) error {
	c, err := wgctrl.New()
	if err != nil {
		return errors.Errorf("failed to open wgctrl: %v", err)
	}
	defer func(c *wgctrl.Client) {
		_ = c.Close()
	}(c)

	var devices []*wgtypes.Device
	if len(args) > 0 {
		device := args[0]
		d, err := c.Device(device)
		if err != nil {
			return errors.Errorf("failed to get device %q: %v", device, err)
		}

		devices = append(devices, d)
	} else {
		devices, err = c.Devices()
		if err != nil {
			return errors.Errorf("failed to get devices: %v", err)
		}
	}

	output := ""
	for _, d := range devices {
		output += formatDevice(d)

		peers := d.Peers
		sort.SliceStable(peers, func(i, j int) bool {
			p1 := peers[i]
			p2 := peers[j]

			return int(time.Since(p1.LastHandshakeTime).Seconds()) > int(time.Since(p2.LastHandshakeTime).Seconds())
		})

		for _, p := range peers {
			output += formatPeer(p)
		}
	}

	fmt.Print(strings.TrimSuffix(output, "\n"))

	return nil
}

func formatDevice(d *wgtypes.Device) string {
	output := greenBoldColor("interface") + ": " + greenColor(d.Name) + " (" + d.Type.String() + ")" + "\n"

	var emptyKey [wgtypes.KeyLen]byte
	if !bytes.Equal(d.PublicKey[:], emptyKey[:]) {
		output +=
			"  " + boldColor("public key") + ": " + d.PublicKey.String() + "\n"
	}
	if !bytes.Equal(d.PrivateKey[:], emptyKey[:]) {
		output +=
			"  " + boldColor("private key") + ": (hidden)\n"
	}

	output +=
		"  " + boldColor("listening port") + ": " + strconv.Itoa(d.ListenPort) + "\n\n"

	return output
}

func formatPeer(p wgtypes.Peer) string {
	output :=
		yellowBoldColor("peer") + ": " + yellowColor(p.PublicKey.String()) + "\n"

	if p.Endpoint != nil {
		output +=
			"  " + boldColor("endpoint") + ": " + p.Endpoint.String() + "\n"
	}

	output +=
		"  " + boldColor("allowed ips") + ": " + strings.ReplaceAll(ipsString(p.AllowedIPs), "/", cyanColor("/")) + "\n"

	if !p.LastHandshakeTime.IsZero() {
		since := int(time.Since(p.LastHandshakeTime).Seconds())
		output +=
			"  " + boldColor("latest handshake") + ": " + formatDuration(since) + " ago" + "\n"
	}

	if p.ReceiveBytes > 0 && p.TransmitBytes > 0 {
		output +=
			"  " + boldColor("transfer") + ": " + formatBytes(p.ReceiveBytes) + " received, " + formatBytes(p.TransmitBytes) + " sent\n"
	}

	output +=
		"  " + boldColor("persistent keepalive") + ": every " + formatDuration(int(p.PersistentKeepaliveInterval.Seconds())) + "\n\n"

	return output
}

func ipsString(ipns []net.IPNet) string {
	ss := make([]string, 0, len(ipns))
	for _, ipn := range ipns {
		ss = append(ss, ipn.String())
	}

	return strings.Join(ss, ", ")
}

func formatDuration(seconds int) string {
	output := ""
	hoursInSeconds := 3600
	dayInSeconds := 24 * hoursInSeconds

	days := seconds / dayInSeconds
	if days > 0 {
		output += formatTimeUnit(days, "day") + ", "
	}

	hours := seconds % dayInSeconds / hoursInSeconds
	if hours > 0 {
		output += formatTimeUnit(hours, "hour") + ", "
	}

	minutes := seconds / 60 % 60
	if minutes > 0 {
		output += formatTimeUnit(minutes, "minute") + ", "
	}

	output += formatTimeUnit(seconds%60, "second")

	return output
}

func formatTimeUnit(value int, unit string) string {
	return strconv.Itoa(value) + " " + plural(value, unit)
}

func plural(value int, unit string) string {
	if value > 1 {
		unit += "s"
	}

	return cyanColor(unit)
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d %s", b, cyanColor("B"))
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.2f %s", float64(b)/float64(div), cyanColor(string("KMGTPE"[exp])+"iB"))
}

func greenBoldColor(s string) string {
	return color.New(color.FgGreen, color.Bold).Sprintf(s)
}

func greenColor(s string) string {
	return color.New(color.FgGreen).Sprintf(s)
}

func boldColor(s string) string {
	return color.New(color.Bold).Sprintf(s)
}

func yellowBoldColor(s string) string {
	return color.New(color.FgYellow, color.Bold).Sprintf(s)
}

func yellowColor(s string) string {
	return color.New(color.FgYellow).Sprintf(s)
}

func cyanColor(s string) string {
	return color.New(color.FgCyan).Sprintf(s)
}
