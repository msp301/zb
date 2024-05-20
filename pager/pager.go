package pager

import (
	"os"

	"github.com/spf13/viper"
)

func FindPager() string {
	pager := viper.GetString("pager")
	if pager != "" {
		return pager
	}

	pager = os.Getenv("PAGER")
	if pager != "" {
		return pager
	}

	return ""
}
