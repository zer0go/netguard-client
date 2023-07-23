package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zer0go/netguard-client/internal/config"
	"github.com/zer0go/netguard-client/internal/handler"
)

var installCmd = &cobra.Command{
	Use:          "install",
	Short:        "Install ngclient binary and daemon",
	SilenceUsage: true,
	RunE:         handler.NewInstallHandler().Handle,
}

func init() {
	installCmd.Flags().StringP("interface", "i", config.DefaultInterfaceName, "interface name")
	installCmd.Flags().StringP("network", "n", "", "network range (eg: 10.2.3.1/24) [required]")
	installCmd.Flags().IntP("mtu", "m", config.DefaultMTU, "mtu")
	_ = installCmd.MarkFlagRequired("network")

	rootCmd.AddCommand(installCmd)
}
