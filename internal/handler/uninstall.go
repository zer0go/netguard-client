package handler

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/zer0go/netguard-client/internal/config"
	"github.com/zer0go/netguard-client/internal/network"
	"os"
)

type UnInstallHandler struct {
}

func NewUnInstallHandler() *UnInstallHandler {
	return &UnInstallHandler{}
}

func (h *UnInstallHandler) Handle(_ *cobra.Command, _ []string) error {
	log.Info().Msg("uninstalling application...")

	if config.IsEmpty() {
		return errors.New("application is not installed")
	}

	networkInterface := network.NewInterfaceFromConfig(config.Get())
	err := networkInterface.LoadLink()
	if err != nil {
		return err
	}

	err = networkInterface.Close()
	if err != nil {
		return err
	}

	err = os.RemoveAll(config.Path)
	if err != nil {
		return err
	}

	log.Info().Msg("uninstalled successfully")

	return nil
}
