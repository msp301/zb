package cmd

import (
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Find anything directly related to a given tag",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		renderResults(book().SearchByTag(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
}
