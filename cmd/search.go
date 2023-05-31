package cmd

import (
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search notes",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		renderResults(book().Search(args...))
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
