package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/msp301/zb/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Get and set configuration options",
	Run: func(cmd *cobra.Command, args []string) {
		configFile := viper.ConfigFileUsed()
		isGlobal, _ := cmd.Flags().GetBool("global")

		if configFile != "" {
			fmt.Printf("Config file used: %s\n\n", viper.ConfigFileUsed())
		} else if isGlobal {
			configFile = config.GlobalConfigFile
		} else {
			configFile = config.ConfigFile
		}

		jsonStr, _ := json.MarshalIndent(viper.AllSettings(), "", "  ")
		fmt.Println(string(jsonStr))

		save, _ := cmd.Flags().GetBool("save")
		if save {
			viper.WriteConfigAs(configFile)
			fmt.Println("Configuration saved to", configFile)
		}
	},
}

func init() {
	configCmd.PersistentFlags().BoolP("global", "g", false, "Use global configuration file")
	configCmd.PersistentFlags().BoolP("save", "s", false, "Write current configuration to file")
	rootCmd.AddCommand(configCmd)
}
