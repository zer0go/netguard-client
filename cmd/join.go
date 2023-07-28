package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zer0go/netguard-client/internal/handler"
)

var joinCommand = &cobra.Command{
	Use:          "join",
	Short:        "Join to server",
	SilenceUsage: true,
	RunE:         handler.NewJoinHandler().Handle,
}

func init() {
	joinCommand.Flags().StringP("token", "t", "", "token")
	_ = joinCommand.MarkFlagRequired("token")

	rootCmd.AddCommand(joinCommand)
}
