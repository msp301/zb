package config

import (
	"os"
	"path/filepath"
	"strings"
)

const CONFIG_NAME = "zb"
const CONFIG_TYPE = "toml"

var ConfigFile string
var GlobalConfigDir string
var GlobalConfigFile string

func init() {
	ConfigFile = strings.Join([]string{CONFIG_NAME, CONFIG_TYPE}, ".")

	GlobalConfigDir = ""
	userConfigDir, err := os.UserConfigDir()
	if err == nil {
		GlobalConfigDir = filepath.Join(userConfigDir, CONFIG_NAME)
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
