package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"path"
)

var exportCmd = &cobra.Command{
	Use:   "export [title]",
	Short: "Export a snippet to json file",
	Args:  cobra.MaximumNArgs(1),
	RunE:  export,
}

var outputFile string

func export(cmd *cobra.Command, args []string) error {
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
	if outputFile == "" {
		outputFile = fmt.Sprintf("./%s", path.Base(s.GetFilePath()))
	}
	if err = s.Export(outputFile); err != nil {
		return err
	}
	return nil
}

func init() {
	exportCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Specify the output path of the snippet")
	rootCmd.AddCommand(exportCmd)
}
