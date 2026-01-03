package cmd

import (
	"runtime"

	"github.com/spf13/cobra"
	"github.com/zer0go/netguard-client/internal/config"
	"github.com/zer0go/netguard-client/internal/handler"
)

var installCmd = &cobra.Command{
	Use:          "install",
	Short:        "Install the application as a system service",
	SilenceUsage: true,
	RunE:         handler.NewInstallHandler().Handle,
}

func init() {
	defaultInterfaceName := config.DefaultInterfaceName
	if runtime.GOOS == "darwin" {
		defaultInterfaceName = config.DefaultDarwinInterfaceName
	}
	installCmd.Flags().StringP("interface", "i", defaultInterfaceName, "interface name")
	installCmd.Flags().Int("mtu", config.DefaultMTU, "mtu")

	rootCmd.AddCommand(installCmd)
}
