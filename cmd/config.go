package cmd

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/msp301/zb/config"
	"github.com/msp301/zb/editor"
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

		if edit, _ := cmd.Flags().GetBool("edit"); edit {
			editor.Open(configFile)
			return
		}

		if len(args) == 0 {
			jsonStr, _ := json.MarshalIndent(viper.AllSettings(), "", "  ")
			fmt.Println(string(jsonStr))

			save, _ := cmd.Flags().GetBool("save")
			if save {
				viper.WriteConfigAs(configFile)
				fmt.Println("Configuration saved to", configFile)
			}

			return
		}

		option := args[0]
		if !viper.InConfig(option) {
			fmt.Println("Option not found")
			return
		}

		if len(args) == 1 {
			value := viper.Get(option)
			fmt.Println(value)
		} else {
			value := viper.Get(option)
			rt := reflect.TypeOf(value)
			switch rt.Kind() {
			case reflect.Array:
			case reflect.Slice:
				viper.Set(option, args[1:])
			default:
				viper.Set(option, args[1])
			}

			viper.WriteConfigAs(configFile)
			fmt.Println("Configuration updated")
		}
	},
}

func init() {
	configCmd.PersistentFlags().BoolP("edit", "e", false, "Open configuration file in editor")
	configCmd.PersistentFlags().BoolP("global", "g", false, "Use global configuration file")
	configCmd.PersistentFlags().BoolP("save", "s", false, "Write current configuration to file")
	rootCmd.AddCommand(configCmd)
}
