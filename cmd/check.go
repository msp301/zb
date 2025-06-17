package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/msp301/zb/util"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check [FILE]",
	Short: "Validate notes and output any that do not pass validation",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			file := args[0]
			fileId, err := util.FileId(file)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			book := book()
			value := []byte{}
			status := 0
			if book.IsNote(fileId) {
				value, _ = json.Marshal(book.Notes.Vertices[fileId])
			} else {
				value, _ = json.Marshal(book.Invalid[fileId])
				status = 1
			}
			fmt.Println(string(value))
			os.Exit(status)
		}

		for _, note := range book().Invalid {
			str, _ := json.Marshal(note)
			fmt.Println(string(str))
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
