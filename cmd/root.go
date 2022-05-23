package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/msp301/zb/graph"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "zb",
	Short: "A Zettelkasten notebook utility",
	Run: func(cmd *cobra.Command, args []string) {
		book().Notes.Walk(func(vertex graph.Vertex, depth int) bool {
			jsonStr, err := json.Marshal(vertex)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(jsonStr))
			return true
		})
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
