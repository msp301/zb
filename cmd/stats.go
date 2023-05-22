package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Output statistics about the notebook",
	Run: func(cmd *cobra.Command, args []string) {
		book := book()
		fmt.Printf("Notes: %d\n", book.Count())
		fmt.Printf("Tags: %d\n", len(book.Tags("")))
		fmt.Printf("Links: %d\n", book.LinkCount())
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
