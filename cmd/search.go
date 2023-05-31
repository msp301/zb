package cmd

import (
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search notes",
	Run: func(cmd *cobra.Command, args []string) {
		renderResults(book().Search(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
