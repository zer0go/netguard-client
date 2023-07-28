package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zer0go/netguard-client/internal/handler"
)

var unInstallCmd = &cobra.Command{
	Use:          "uninstall",
	Short:        "Uninstall ngclient binary and daemon",
	SilenceUsage: true,
	RunE:         handler.NewUnInstallHandler().Handle,
}

func init() {
	rootCmd.AddCommand(unInstallCmd)
}
