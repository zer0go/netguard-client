package handler

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/zer0go/netguard-client/internal/config"
	"github.com/zer0go/netguard-client/internal/network"
	"os"
)

type InstallHandler struct {
}

func NewInstallHandler() *InstallHandler {
	return &InstallHandler{}
}

func (h *InstallHandler) Handle(cmd *cobra.Command, _ []string) error {
	log.Debug().Msg("installing application...")

	err := install()
	if err != nil {
		return err
	}

	interfaceName, _ := cmd.Flags().GetString("interface")
	mtu, _ := cmd.Flags().GetInt("mtu")

	if _, err := os.Stat(config.GetConfigPath()); os.IsNotExist(err) {
		if err := os.Mkdir(config.GetConfigPath(), os.ModePerm); err != nil {
			return err
		}
	}

	err = config.Update(config.App{
		InterfaceName: interfaceName,
		MTU:           mtu,
	})
	if err != nil {
		return err
	}

	appConfig := config.Get()
	log.Debug().
		Interface("AppConfig", appConfig).
		Msg("config updated")

	networkInterface := network.NewInterfaceFromConfig(appConfig)
	if err := networkInterface.Create(); err != nil {
		return err
	}
	if err := networkInterface.Configure(); err != nil {
		return err
	}

	log.Debug().Msg("application installed.")

	return nil
}
