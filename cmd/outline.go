package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
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
		pager := pager.FindPager()
		if pager == "" {
			log.Fatal("No pager found")
		}

		command := exec.Command(pager)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		pipe, err := command.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			defer pipe.Close()

			book().Notes.Walk(func(vertex graph.Vertex, depth int) bool {
				indent := strings.Repeat("\t", depth)
				switch val := vertex.Properties["Value"].(type) {
				case parser.Note:
					_, err := io.WriteString(pipe, fmt.Sprintf("%s%s - %s\n", indent, val.File, val.Title))
					if err != nil {
						log.Fatal(err)
					}
				}
				return true
			}, -1)
		}()

		err = command.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(outlineCmd)
}
