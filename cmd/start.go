package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zer0go/netguard-client/internal/handler"
)

var startCmd = &cobra.Command{
	Use:          "start",
	Short:        "Start the VPN connection(s)",
	SilenceUsage: true,
	RunE:         handler.NewStartHandler().Handle,
}

func init() {
	rootCmd.AddCommand(startCmd)
}
