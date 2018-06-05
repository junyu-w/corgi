package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a snippet",
	Args:  cobra.ExactArgs(1),
	RunE:  edit,
}

func edit(cmd *cobra.Command, args []string) error {
	title := args[0]
	// load config & snippets
	config, snippets, err := loadConfigAndSnippetsMeta()
	if err != nil {
		return err
	}
	// find snippet
	s, err := snippets.FindSnippet(title)
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
	rootCmd.AddCommand(editCmd)
}
