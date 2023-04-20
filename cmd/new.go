package cmd

import (
	"fmt"
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
	Run: func(cmd *cobra.Command, args []string) {
		dirs := bookDirs()

		notePath := filepath.Join(dirs[0], strconv.Itoa(time.Now().Year()))
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

		fmt.Println(fullPath)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func createNote(dir string) (*os.File, error) {
	timestamp := time.Now().Format("200601021504")

	note, err := os.Create(filepath.Join(dir, timestamp+".md"))
	if err != nil {
		return nil, err
	}
	defer note.Close()

	note.WriteString(timestamp)

	return note, nil
}
