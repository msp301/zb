package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate notes and output any that do not pass validation",
	Run: func(cmd *cobra.Command, args []string) {
		for _, note := range book().Invalid {
			str, _ := json.Marshal(note)
			fmt.Println(string(str))
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
