package cmd

import (
	"github.com/msp301/zb/notebook"
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Find anything directly related to one or more tags",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var results []notebook.Result
		if len(args) == 1 {
			results = book().SearchByTag(args[0])
		} else {
			results = book().SearchByTags(args)
		}
		renderResults(results)
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
}
