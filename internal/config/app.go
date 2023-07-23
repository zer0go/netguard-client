package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

const (
	Path     = "/etc/netguard"
	FileName = "config.yaml"

	DefaultInterfaceName = "netguard"
	DefaultMTU           = 1280
	DefaultListeningPort = 51821
)

var (
	config App
)

type App struct {
	InterfaceName string `yaml:"interface_name"`
	MTU           int    `yaml:"mtu"`
	NetworkCIDR   string `yaml:"network_cidr"`
}

func init() {
	config = App{}
}

func Load() error {
	data, err := os.ReadFile(Path + "/" + FileName)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	return nil
}

func Update(c App) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(Path+"/"+FileName, data, 0644)
	if err != nil {
		return err
	}

	config = c

	return nil
}

func Get() *App {
	return &config
}