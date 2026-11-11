package cmd

import (
	"errors"
	"fmt"
	"github.com/msp301/zb/editor"
	"github.com/spf13/viper"
	"io/fs"
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

		fullPath, err := filepath.Abs(note)
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

func createNote(dir string) (string, error) {
	for created := time.Now(); ; created = created.Add(time.Minute) {
		timestamp := created.Format("200601021504")
		path := filepath.Join(dir, timestamp+".md")

		note, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
		if errors.Is(err, fs.ErrExist) {
			continue
		}
		if err != nil {
			return "", err
		}

		if _, err = note.WriteString(timestamp); err != nil {
			note.Close()
			return "", err
		}

		return path, note.Close()
	}
}

func aliasIndex(alias string) int {
	for i, a := range viper.GetStringSlice("alias") {
		if a == alias {
			return i
		}
	}
	return -1
}
