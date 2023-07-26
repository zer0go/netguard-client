package cmd

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/cobra"
	"github.com/zer0go/netguard-client/internal/config"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "ngclient",
	Short: "NetGuard Client",
}

func init() {
	var verbosity int
	rootCmd.
		PersistentFlags().
		IntVarP(&verbosity, "verbosity", "v", 0, "set logging verbosity 0-4")
	_ = rootCmd.ParseFlags(os.Args[1:])

	switch verbosity {
	case 4:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case 3:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case 2:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}
}

func Execute(version string) {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		dir := filepath.Dir(file)
		parent := filepath.Base(dir)
		return parent + "/" + filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	log.Logger = log.
		With().
		Stack().
		Caller().
		Str("app_version", version).
		Logger()

	err := config.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Warn().Err(err).Msg("")
	}

	rootCmd.Version = version
	rootCmd.Short += " " + version

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err)
	}
}
