package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec [title]",
	Short: "Execute a snippet",
	Args:  cobra.MaximumNArgs(1),
	RunE:  execute,
}

var useDefaultParamValue bool
var stepRange string

func execute(cmd *cobra.Command, args []string) error {
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
	// find snippet corresponds to title
	s, err := snippetsMeta.FindSnippet(title)
	if err != nil {
		return fmt.Errorf("%s, run \"corgi list\" to view all snippets", err.Error())
	}
	s.Execute(useDefaultParamValue, stepRange)
	return nil
}

func init() {
	execCmd.Flags().StringVarP(&stepRange, "step", "s", "", "Select a single step to execute with \"-s <step>\" or a range of steps to execute with \"-s <start>-<end>\", end is optional")
	execCmd.Flags().BoolVar(&useDefaultParamValue, "use-default", false, "Add this flag if you would like to use the default values for your defined template fields without being asked to enter a value")
	rootCmd.AddCommand(execCmd)
}
