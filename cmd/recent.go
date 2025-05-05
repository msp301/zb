package cmd

import (
	"fmt"
	"github.com/msp301/graph"
	"github.com/msp301/zb/parser"
	"github.com/spf13/cobra"
)

var recentCmd = &cobra.Command{
	Use:   "recent",
	Short: "Output notes recently worked on",
	Run: func(cmd *cobra.Command, args []string) {
		maxResults, _ := cmd.Flags().GetInt("number")
		results := 0
		book().Notes.WalkBackwards(func(vertex graph.Vertex, _ int) bool {
			if results >= maxResults {
				return false
			}
			switch val := vertex.Properties["Value"].(type) {
			case parser.Note:
				fmt.Printf("%s - %s\n", val.File, val.Title)
				results++
			}
			return true
		}, 0)
	},
}

func init() {
	recentCmd.PersistentFlags().IntP("number", "n", 10, "Number of results to return")
	rootCmd.AddCommand(recentCmd)
}
