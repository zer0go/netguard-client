package config

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	AppName              = "NetGuard"
	DefaultInterfaceName = "netguard"
	DefaultMTU           = 1280
	FileName             = "config.yaml"
	WireGuardNTUrl       = "https://download.wireguard.com/wireguard-nt/wireguard-nt-0.10.1.zip"
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

func Get() *App {
	return &config
}

func Load() error {
	data, err := os.ReadFile(filepath.Join(GetConfigPath(), FileName))
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

	err = os.WriteFile(filepath.Join(GetConfigPath(), FileName), data, os.ModePerm)
	if err != nil {
		return err
	}

	config = c

	return nil
}

func IsEmpty() bool {
	stat, err := os.Stat(filepath.Join(GetConfigPath(), FileName))
	if err != nil {
		return true
	}

	return stat.Size() == 0
}

func ConfigureLogger(verbosity int) {
	switch verbosity {
	case 1:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case 0:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		dir := filepath.Dir(file)
		parent := filepath.Base(dir)
		return parent + "/" + filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	log.Logger = log.
		With().
		Timestamp().
		Stack().
		Caller().
		Logger()

	if os.Getenv("LOG_FORMAT") != "json" {
		log.Logger = log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: zerolog.TimeFieldFormat})
	}
}
