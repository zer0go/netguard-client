package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zer0go/netguard-client/internal/network"
	"strconv"
)

var stunCommand = &cobra.Command{
	Use:          "stun <port>",
	Short:        "UDP hole punching on the given port",
	SilenceUsage: true,
	Args:         cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		portToStun, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		i, p, t, err := network.NewStunClientWithServers().StunPort(portToStun)
		if err != nil {
			return err
		}

		fmt.Println(i, p, t)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(stunCommand)
}
