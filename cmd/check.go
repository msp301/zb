package cmd

import (
	"github.com/msp301/zb/notebook"
	"github.com/msp301/zb/parser"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate notes and output any that do not pass validation",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Fix 'check' command to return invalid notes
		book := book()
		book.AddFilter(func(note parser.Note) bool {
			return isValidNote(note, book)
		})
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func isValidNote(note parser.Note, book *notebook.Notebook) bool {
	if note.Id == 0 {
		return false
	}

	for _, link := range note.Links {
		if !book.IsNote(link) {
			return false
		}
	}

	return true
}
