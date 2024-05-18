package cmd

import (
	"fmt"
	"github.com/msp301/zb/editor"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new note",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dirs := bookDirs()
		bookDir := dirs[0]

		if len(args) > 0 {
			bookAlias := args[0]

			aliasIndex := aliasIndex(bookAlias)
			if aliasIndex != -1 {
				bookDir = dirs[aliasIndex]
			}
		}

		notePath := filepath.Join(bookDir, strconv.Itoa(time.Now().Year()))
		if err := os.MkdirAll(notePath, 0755); err != nil {
			log.Fatalf("Error creating directory: %s", err)
		}

		note, err := createNote(notePath)
		if err != nil {
			log.Fatalf("Error creating note: %s", err)
		}

		fullPath, err := filepath.Abs(note.Name())
		if err != nil {
			log.Fatalf("Error getting absolute path: %s", err)
		}

		err = editor.Open(fullPath)
		if err != nil {
			log.Printf("Editor error: %s\n", err)
			fmt.Println(fullPath)
			os.Exit(1)
		}
	},
}

func init() {
	newCmd.PersistentFlags().StringSlice("alias", []string{}, "aliases for notebook directories")
	viper.BindPFlag("alias", newCmd.PersistentFlags().Lookup("alias"))

	rootCmd.AddCommand(newCmd)
}

func createNote(dir string) (*os.File, error) {
	timestamp := time.Now().Format("200601021504")

	note, err := os.Create(filepath.Join(dir, timestamp+".md"))
	if err != nil {
		return nil, err
	}
	defer func(note *os.File) {
		err := note.Close()
		if err != nil {
			log.Fatalf("Error closing note: %s", err)
		}
	}(note)

	_, err = note.WriteString(timestamp)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func aliasIndex(alias string) int {
	for i, a := range viper.GetStringSlice("alias") {
		if a == alias {
			return i
		}
	}
	return -1
}
