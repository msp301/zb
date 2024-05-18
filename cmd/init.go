package cmd

import (
	"log"
	"path/filepath"

	"github.com/msp301/zb/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Add current directory to notebook",
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

		viper.Set("directory", append(viper.GetStringSlice("directory"), fullPath))

		if len(args) > 0 {
			alias := args[0]
			viper.Set("alias", append(viper.GetStringSlice("alias"), alias))
		}

		writeConfig(configFile)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
