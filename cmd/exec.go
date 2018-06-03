package cmd

import (
	"github.com/kataras/iris/core/errors"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a snippet",
	RunE:  exec,
}

var execTitle string

func exec(cmd *cobra.Command, args []string) error {
	if execTitle == "" {
		// TODO: launch fzf search
		return errors.New("must specify --title to execute command")
	}
	// load config & snippets
	_, snippets, err := loadConfigAndSnippets()
	if err != nil {
		return err
	}
	// find snippet corresponds to title
	s, err := snippets.FindSnippet(execTitle)
	if err != nil {
		return err
	}
	if err = s.Execute(); err != nil {
		return err
	}
	return nil
}

func init() {
	execCmd.Flags().StringVarP(&execTitle, "title", "t", "", "Name of the snippet to execute")
	rootCmd.AddCommand(execCmd)
}
