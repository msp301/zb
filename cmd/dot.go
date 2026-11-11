package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var dotCmd = &cobra.Command{
	Use:   "dot",
	Short: "Output a DOT graph of the subgraph matched by a tag query",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			args = strings.Fields(args[0])
		}
		fmt.Print(book().VisualizeTagQuery(args...))
	},
}

func init() {
	rootCmd.AddCommand(dotCmd)
}
