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
	installCmd.Flags().Int("mtu", config.DefaultMTU, "mtu")

	rootCmd.AddCommand(installCmd)
}
