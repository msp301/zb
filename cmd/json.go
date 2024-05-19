package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/msp301/zb/graph"
	"github.com/spf13/cobra"
)

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "Output notes as JSON",
	Run: func(cmd *cobra.Command, args []string) {
		book().Notes.Walk(func(vertex graph.Vertex, depth int) bool {
			jsonStr, err := json.Marshal(vertex)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(jsonStr))
			return true
		}, -1)
	},
}

func init() {
	rootCmd.AddCommand(jsonCmd)
}
