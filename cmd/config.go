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
		jsonStr, _ := json.MarshalIndent(viper.AllSettings(), "", "  ")
		fmt.Println(string(jsonStr))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
