package cmd

import (
	"log"

	"github.com/msp301/zb/pager"
	"github.com/spf13/cobra"
)

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List tags / Search for a tag by given string",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		searchTag := ""
		if len(args) >= 1 {
			searchTag = args[0]
		}

		pager, err := pager.Open()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			defer pager.Close()
			if connections, err := cmd.Flags().GetBool("connections"); err == nil && connections {
				for _, tagConnection := range book().TagConnections(searchTag) {
					pager.Writef("%d %s\n", tagConnection.Connections, tagConnection.Tag)
				}
				return
			}

			for _, tag := range book().Tags(searchTag) {
				pager.Writeln(tag)
			}
		}()

		if err := pager.Wait(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tagsCmd)

	tagsCmd.Flags().BoolP("connections", "c", false, "Include number of connections to tags")
}
