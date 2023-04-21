package cmd

import (
	"fmt"
	"github.com/msp301/zb/graph"
	"github.com/msp301/zb/parser"
	"github.com/spf13/cobra"
	"math/rand"
)

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "List random selection of notes to review",
	Run: func(cmd *cobra.Command, args []string) {
		var notes []graph.Vertex
		for _, vertex := range book().Notes.Vertices {
			notes = append(notes, vertex)
		}
		shuffledNotes := shuffle(notes)
		shuffledNotesLen := len(shuffledNotes)
		num := 5
		if shuffledNotesLen < num {
			num = shuffledNotesLen
		}
		count := 0
		for _, vertex := range shuffledNotes {
			if count > num {
				break
			}
			switch val := vertex.Properties["Value"].(type) {
			case parser.Note:
				fmt.Printf("%s - %s\n", val.File, val.Title)
				count++
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(reviewCmd)
}

func shuffle(arr []graph.Vertex) []graph.Vertex {
	for i := range arr {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}
