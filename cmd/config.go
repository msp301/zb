package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Get and set configuration options",
	Run: func(cmd *cobra.Command, args []string) {
		configFile := viper.ConfigFileUsed()
		if configFile != "" {
			fmt.Printf("Config file used: %s\n\n", viper.ConfigFileUsed())
		}

		jsonStr, _ := json.MarshalIndent(viper.AllSettings(), "", "  ")
		fmt.Println(string(jsonStr))

		save, _ := cmd.Flags().GetBool("save")
		if save {
			viper.WriteConfig()
			fmt.Println("Configuration saved to", viper.ConfigFileUsed())
		}
	},
}

func init() {
	configCmd.PersistentFlags().BoolP("save", "s", false, "Write current configuration to file")
	viper.BindPFlag("save", configCmd.PersistentFlags().Lookup("save"))
	rootCmd.AddCommand(configCmd)
}
