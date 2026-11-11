package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"slices"

	"github.com/msp301/zb/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init [ALIAS]",
	Short: "Add current directory to notebook, optionally under an alias",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		directory := defaultNotebookDir()

		fullPath, err := filepath.Abs(directory)
		if err != nil {
			log.Fatalf("Error getting absolute path: %s", err)
		}

		configFile := viper.ConfigFileUsed()
		if configFile == "" {
			configFile = config.GlobalConfigFile
		}

		dirs := viper.GetStringSlice("directory")
		index := slices.Index(dirs, fullPath)
		if index == -1 {
			dirs = append(dirs, fullPath)
			index = len(dirs) - 1
		} else if len(args) == 0 {
			fmt.Printf("%s is already part of the notebook\n", fullPath)
			return
		}
		viper.Set("directory", dirs)

		if len(args) > 0 {
			alias := args[0]
			if existing := aliasIndex(alias); existing != -1 && existing != index {
				fail("Alias '%s' is already in use", alias)
			}

			aliases := viper.GetStringSlice("alias")
			for len(aliases) <= index {
				aliases = append(aliases, "")
			}
			aliases[index] = alias
			viper.Set("alias", aliases)
		}

		writeConfig(configFile)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
