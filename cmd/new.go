package cmd

import (
	"corgi/snippet"
	"fmt"
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
	fmt.Println("New command called")
	hcmds, err := snippet.ReadShellHistory(lastCmds)
	if err != nil {
		return err
	}
	fmt.Println(hcmds)
	newSnippet, err := snippet.NewSnippet(hcmds)
	if err != nil {
		return err
	}
	if err := newSnippet.Save(); err != nil {
		return err
	}
	return nil
}

func init() {
	newCmd.Flags().IntVarP(&lastCmds, "last", "l", 1, "Select the number of history commands to look back")
	newCmd.Flags().BoolVar(&withoutDescription, "without-description", false, "Skip entering description (use command itself as default)")
	newCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the snippet")
	rootCmd.AddCommand(newCmd)
}
