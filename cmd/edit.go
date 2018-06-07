package cmd

import (
	"corgi/util"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a snippet",
	Args:  cobra.MaximumNArgs(1),
	RunE:  edit,
}

func edit(cmd *cobra.Command, args []string) error {
	// load config & snippets
	conf, snippets, err := loadConfigAndSnippetsMeta()
	if err != nil {
		return err
	}
	// find snippet title
	var title string
	if len(args) == 0 {
		if conf.FilterCmd != "" {
			title, err = filter(conf.FilterCmd, snippets.GetSnippetTitles())
			if err != nil || title == "" {
				return MissingSnippetTitleError
			}
		} else {
			color.Red("Install a fuzzy finder (\"fzf\" or \"peco\") to enable interactive selection")
			return MissingSnippetTitleError
		}
	} else {
		title = args[0]
	}
	// find snippet
	s, err := snippets.FindSnippet(title)
	if err != nil {
		return err
	}
	command := fmt.Sprintf("%s %s", conf.Editor, s.GetFilePath())
	if err := util.Execute(command, os.Stdin, os.Stdout); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(editCmd)
}
