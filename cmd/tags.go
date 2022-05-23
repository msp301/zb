package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List tags / Search for a tag by given string",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		searchTag := ""
		if len(args) >= 1 {
			searchTag = args[0]
		}

		for _, tag := range book().Tags(searchTag) {
			fmt.Println(tag)
		}
	},
}

func init() {
	rootCmd.AddCommand(tagsCmd)
}
