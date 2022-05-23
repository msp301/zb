package cmd

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Find anything directly related to a given note",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := strconv.ParseUint(os.Args[3], 0, 64)
		render(book().SearchRelated(id))
	},
}

func init() {
	rootCmd.AddCommand(noteCmd)
}
