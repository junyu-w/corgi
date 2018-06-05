package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a snippet",
	Args:  cobra.ExactArgs(1),
	RunE:  execute,
}

func execute(cmd *cobra.Command, args []string) error {
	title := args[0]
	// load config & snippets
	_, snippets, err := loadConfigAndSnippetsMeta()
	if err != nil {
		return err
	}
	// find snippet corresponds to title
	s, err := snippets.FindSnippet(title)
	if err != nil {
		return fmt.Errorf("%s, run \"corgi list\" to view all snippets", err.Error())
	}
	s.Execute()
	return nil
}

func init() {
	rootCmd.AddCommand(execCmd)
}
