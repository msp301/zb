package cmd

import (
	"fmt"
	"strconv"

	"github.com/msp301/zb/util"
	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Find anything directly related to a given note",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		note := args[0]
		id, err := strconv.ParseUint(note, 0, 64)
		if err != nil {
			id, err = util.FileId(note)
			if err != nil {
				fmt.Println(err)
			}
		}
		renderResults(book().SearchRelated(id))
	},
}

func init() {
	rootCmd.AddCommand(noteCmd)
}
