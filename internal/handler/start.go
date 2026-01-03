package handler

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type StartHandler struct {
}

func NewStartHandler() *StartHandler {
	return &StartHandler{}
}

func (h *StartHandler) Handle(_ *cobra.Command, _ []string) error {
	log.Info().Msg("starting application...")
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	err, closeFunc := startNetworkAdapter()
	if err != nil {
		return err
	}

	log.Info().Msg("application started successfully.")

	<-ctx.Done()

	closeFunc()

	log.Info().Msg("shutting down application.")

	return nil
}
