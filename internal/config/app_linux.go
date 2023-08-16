package config

import (
	"path/filepath"
	"strings"
)

func GetConfigPath() string {
	return filepath.Join("/", "etc", strings.ToLower(AppName))
}
