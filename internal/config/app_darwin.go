package config

import (
	"path/filepath"
	"strings"
)

func GetConfigPath() string {
	return filepath.Join("/", "etc", strings.ToLower(AppName))
}

func GetLogPath() string {
	return filepath.Join("/", "var", "log")
}
