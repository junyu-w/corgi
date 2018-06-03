package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a snippet",
	RunE:  edit,
}

var editTitle string

func edit(cmd *cobra.Command, args []string) error {
	if editTitle == "" {
		// TODO: launch fzf search
		return errors.New("must specify --title to edit command")
	}
	// load config & snippets
	config, snippets, err := loadConfigAndSnippets()
	if err != nil {
		return err
	}
	// find snippet
	s, err := snippets.FindSnippet(editTitle)
	if err != nil {
		return err
	}
	editFileCmd := exec.Command(config.Editor, s.GetFilePath())
	editFileCmd.Stdin = os.Stdin
	editFileCmd.Stdout = os.Stdout
	if err := editFileCmd.Run(); err != nil {
		return err
	}
	return nil
}

func init() {
	editCmd.Flags().StringVarP(&editTitle, "title", "t", "", "Name of the snippet to edit with $EDITOR, vim (default), or an editor of your choice (set by the configure command)")
	rootCmd.AddCommand(editCmd)
}
