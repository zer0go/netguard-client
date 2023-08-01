package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/zer0go/netguard-client/internal/config"
	"os"
)

var rootCmd = &cobra.Command{
	Use:              "ngclient",
	Short:            "NetGuard Client",
	SilenceErrors:    true,
	PersistentPreRun: bootstrap,
}

func init() {
	rootCmd.
		PersistentFlags().
		CountP("verbosity", "v", "set logging verbosity")
}

func bootstrap(cmd *cobra.Command, _ []string) {
	verbosity, _ := cmd.Flags().GetCount("verbosity")
	config.ConfigureLogger(verbosity)

	err := config.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Warn().Err(err).Msg("")
	}
}

func Execute(version string) {
	rootCmd.Version = version
	rootCmd.Short += " " + version

	if err := rootCmd.Execute(); err != nil {
		log.Warn().Err(err).Msg("unexpected error")
		os.Exit(1)
	}
}
