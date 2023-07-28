package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zer0go/netguard-client/internal/handler"
)

var showCmd = &cobra.Command{
	Use:          "show { interface }",
	Short:        "Show wireguard config and device information",
	SilenceUsage: true,
	RunE:         handler.NewShowHandler().Handle,
	Version:      "0.1.8",
}

func init() {
	rootCmd.AddCommand(showCmd)
}
