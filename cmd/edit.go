package cmd

import (
	"corgi/util"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var editCmd = &cobra.Command{
	Use:   "edit [title]",
	Short: "Edit a snippet",
	Args:  cobra.MaximumNArgs(1),
	RunE:  edit,
}

func edit(cmd *cobra.Command, args []string) error {
	// load config & snippets
	conf, snippetsMeta, err := loadConfigAndSnippetsMeta()
	if err != nil {
		return err
	}
	// find snippet title
	var title string
	if len(args) == 0 {
		title, err = filterSnippetTitle(conf.FilterCmd, snippetsMeta.GetSnippetTitles())
		if err != nil {
			return err
		}
	} else {
		title = args[0]
	}
	// find snippet
	s, err := snippetsMeta.FindSnippet(title)
	if err != nil {
		return err
	}
	command := fmt.Sprintf("%s %s", conf.Editor, s.GetFilePath())
	if err := util.Execute(command, os.Stdin, os.Stdout); err != nil {
		return err
	}
	// mark snippetsMeta dirty
	snippetsMeta.IsMetaDirty = true
	if err = snippetsMeta.Save(); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(editCmd)
}
