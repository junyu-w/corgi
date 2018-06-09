package cmd

import (
	"corgi/snippet"
	"github.com/spf13/cobra"
)

// flags
var lastCmds int
var title string

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new corgi snippet",
	Long: `Create a new corgi snippet from your command line history interactively

Note: if you plan to use other corgi command that takes a snippet title (for example: corgi exec -t <title>), make sure you don't put double quotes around <title>, otherwise weird failure will happen`,
	RunE: create,
}

func create(cmd *cobra.Command, args []string) error {
	// set up history
	histCmds, err := snippet.ReadShellHistory()
	if err != nil {
		return err
	}
	if err = snippet.SetUpHistFile(histCmds); err != nil {
		return err
	}
	defer snippet.RemoveHistFile()
	// load config and snippets
	conf, snippetsMeta, err := loadConfigAndSnippetsMeta()
	if err != nil {
		return err
	}
	// create snippet
	initialDefaultCmds := histCmds[len(histCmds)-(lastCmds+1) : len(histCmds)-1]
	newSnippet, err := snippet.NewSnippet(title, initialDefaultCmds)
	if err != nil {
		return err
	}
	// add new sninppet to snippets meta and save
	if err = snippetsMeta.SaveNewSnippet(newSnippet, conf.SnippetsDir); err != nil {
		return err
	}
	return nil
}

func init() {
	newCmd.Flags().IntVarP(&lastCmds, "last", "l", 0, "The number of history commands to look back, they'll be the default for each step. If 0 or unspecified, each step will not have a default.")
	newCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the snippet, do not put any whitespace if you plan to use this snippet for composition")
	rootCmd.AddCommand(newCmd)
}
