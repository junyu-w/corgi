package cmd

import (
	"corgi/config"
	"corgi/snippet"
	"github.com/spf13/cobra"
)

// flags
var lastCmds int
var withoutDescription bool
var title string

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new corgi snippet",
	Long:  "Create a new corgi snippet from your command line history interactively",
	RunE:  create,
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
	// Load config
	conf, err := config.Load()
	if err != nil {
		return err
	}
	// Load snippets
	snippets, err := snippet.LoadSnippets(conf.SnippetsFile)
	if err != nil {
		return err
	}
	// create snippet
	initialDefaultCmds := histCmds[len(histCmds)-(lastCmds+1) : len(histCmds)-1]
	newSnippet, err := snippet.NewSnippet(title, initialDefaultCmds)
	if err != nil {
		return err
	}
	if err := newSnippet.Save(conf.SnippetsDir); err != nil {
		return err
	}
	// add new sninppet to snippets and save
	snippets.AddSnippet(newSnippet)
	if err = snippets.Save(); err != nil {
		return err
	}
	return nil
}

func init() {
	newCmd.Flags().IntVarP(&lastCmds, "last", "l", 0, "The number of history commands to look back, they'll be the default for each step. If 0 or unspecified, each step will not have a default.")
	newCmd.Flags().BoolVar(&withoutDescription, "without-description", false, "Skip entering description (use command itself as default)")
	newCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the snippet")
	rootCmd.AddCommand(newCmd)
}
