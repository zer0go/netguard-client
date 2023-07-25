package cmd

import (
	"github.com/spf13/cobra"
	//	"github.com/zer0go/netguard-client/internal/config"
	"github.com/zer0go/netguard-client/internal/handler"
)

var command = &cobra.Command{
	Use:          "join",
	Short:        "join to server",
	SilenceUsage: true,
	RunE:         handler.NewJoinHandler().Handle,
}

func init() {
	command.Flags().StringP("token", "t", "", "token")
	_ = command.MarkFlagRequired("token")

	rootCmd.AddCommand(command)
}
