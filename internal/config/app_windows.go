package config

import (
	"os"
	"path/filepath"
)

func GetConfigPath() string {
	return filepath.Join(os.Getenv("PROGRAMDATA"), AppName)
}
