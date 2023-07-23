package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zer0go/netguard-client/internal/config"
	"github.com/zer0go/netguard-client/internal/handler"
)

var command = &cobra.Command{
	Use:          "join",
	Short:        "join to server",
	SilenceUsage: true,
	RunE:         handler.NewJoinHandler().Handle,
}

func init() {
	command.Flags().StringP("interface", "i", config.DefaultInterfaceName, "interface name")
	command.Flags().String("private_key", "", "private key [required]")
	command.Flags().String("peer_allowed_ip", "", "peer allowed ip [required]")
	command.Flags().String("peer_endpoint", "", "peer endpoint [required]")
	command.Flags().String("peer_public_key", "", "peer public key [required]")
	command.Flags().IntP("listening_port", "p", config.DefaultListeningPort, "listening port")

	_ = command.MarkFlagRequired("private_key")
	_ = command.MarkFlagRequired("peer_allowed_ip")
	_ = command.MarkFlagRequired("peer_endpoint")
	_ = command.MarkFlagRequired("peer_public_key")

	rootCmd.AddCommand(command)
}
