package cmd

import (
	"log"
	"strings"

	"github.com/msp301/zb/graph"
	"github.com/msp301/zb/pager"
	"github.com/msp301/zb/parser"
	"github.com/spf13/cobra"
)

var outlineCmd = &cobra.Command{
	Use:   "outline",
	Short: "Output notes as a tree",
	Run: func(cmd *cobra.Command, args []string) {
		pager, err := pager.Open()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			defer pager.Close()

			book().Notes.Walk(func(vertex graph.Vertex, depth int) bool {
				indent := strings.Repeat("\t", depth)
				switch val := vertex.Properties["Value"].(type) {
				case parser.Note:
					_, err := pager.Writef("%s%s - %s\n", indent, val.File, val.Title)
					if err != nil {
						log.Fatal(err)
					}
				}
				return true
			}, -1)
		}()

		if err := pager.Wait(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(outlineCmd)
}
