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

		if configFile == "" {
			configFile = config.GlobalConfigFile
		}

		if configFile != config.GlobalConfigFile {
			fmt.Printf("Config file used: %s\n\n", viper.ConfigFileUsed())
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
		if !config.IsConfigOption(option) {
			fmt.Printf("Option '%s' not found\n", option)
			return
		}

		if unset, _ := cmd.Flags().GetBool("unset"); unset {
			// Viper does not have in-built support for deleting or unsetting a config key (https://github.com/spf13/viper/pull/519)
			value := viper.Get(option)
			rt := reflect.TypeOf(value)
			switch rt.Kind() {
			case reflect.Array:
			case reflect.Slice:
				viper.Set(option, []string{})
			default:
				viper.Set(option, "")
			}

			viper.WriteConfigAs(configFile)
			fmt.Println("Configuration updated")
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
	configCmd.PersistentFlags().BoolP("save", "s", false, "Write current configuration to file")
	configCmd.PersistentFlags().BoolP("unset", "u", false, "Delete configuration option value")
	rootCmd.AddCommand(configCmd)
}
