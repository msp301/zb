package cmd

import (
	"strings"

	"github.com/msp301/zb/bookshelf"
	"github.com/msp301/zb/browser"
	"github.com/msp301/zb/notebook"
	"github.com/spf13/cobra"
)

var browseCmd = &cobra.Command{
	Use:   "browse [QUERY]...",
	Short: "Interactively search, read and navigate between notes",
	Long: `Interactively search, read and navigate between notes.

Type to search note text, or start terms with '#' to search by tag.
With an empty query the most recent notes are listed.

  Up / Down          select a result
  PgUp / PgDn        scroll the note being read
  Tab                show what is related to the selected note or tag
  Shift+Tab          go back to the previous list
  Enter              open the selected note in your editor
  Esc                clear the query, go back, then quit
  Ctrl+C             quit

To open notes at the matched line, include placeholders in the editor
setting, for example:

  zb config editor "nvim +{line} {file}"
  zb config editor "code --goto {file}:{line}"`,
	Run: func(cmd *cobra.Command, args []string) {
		dirs := bookDirs()
		err := browser.Run(browser.Options{
			Dirs:  dirs,
			Query: strings.Join(args, " "),
			Load: func() (*notebook.Notebook, error) {
				err, book := bookshelf.Read(dirs)
				return book, err
			},
		})
		if err != nil {
			fail("%s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(browseCmd)
}
