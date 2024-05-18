package config

import (
	"os"
	"path/filepath"
	"strings"
)

const CONFIG_NAME = ".zb"
const CONFIG_TYPE = "toml"

var ConfigFile string
var GlobalConfigDir string
var GlobalConfigFile string

func init() {
	ConfigFile = strings.Join([]string{CONFIG_NAME, CONFIG_TYPE}, ".")

	GlobalConfigDir = ""
	home, err := os.UserHomeDir()
	if err == nil {
		GlobalConfigDir = home
	}

	GlobalConfigFile = filepath.Join(GlobalConfigDir, ConfigFile)
}

func IsConfigOption(key string) bool {
	switch key {
	case "alias", "directory", "editor":
		return true
	}

	return false
}
