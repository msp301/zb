package cmd

import (
	"fmt"
	"github.com/msp301/zb/graph"
	"github.com/msp301/zb/parser"
	"strings"

	"github.com/spf13/cobra"
)

var outlineCmd = &cobra.Command{
	Use:   "outline",
	Short: "Output notes as a tree",
	Run: func(cmd *cobra.Command, args []string) {
		book().Notes.Walk(func(vertex graph.Vertex, depth int) bool {
			indent := strings.Repeat("\t", depth)
			switch val := vertex.Properties["Value"].(type) {
			case parser.Note:
				fmt.Printf("%s%s - %s\n", indent, val.File, val.Title)
			}
			return true
		})
	},
}

func init() {
	rootCmd.AddCommand(outlineCmd)
}
