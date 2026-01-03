package config

import (
	"os"
	"path/filepath"
)

func GetConfigPath() string {
	return filepath.Join(os.Getenv("PROGRAMDATA"), AppName)
}

func GetLogPath() string {
	return filepath.Join(GetConfigPath(), "logs")
}
