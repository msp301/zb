package cmd

import (
	"github.com/spf13/cobra"
	"strings"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Find anything directly related to one or more tags",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			args = strings.Fields(args[0])
		}
		renderResults(book().SearchByTags(args...))
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
}
