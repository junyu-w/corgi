package cmd

import (
	"fmt"
	"path"
	"strings"

	"github.com/DrakeW/corgi/snippet"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export [title]",
	Short: "Export a snippet to json file",
	Args:  cobra.MaximumNArgs(1),
	RunE:  export,
}

var outputFile string
var fileType string

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
		corgiFileName := path.Base(s.GetFilePath())
		if fileType == snippet.EXPORT_TYPE_CORGI {
			outputFile = fmt.Sprintf("./%s", corgiFileName)
		} else {
			outputFile = fmt.Sprintf("./%s", strings.Split(corgiFileName, ".")[0])
		}
	}
	err = s.Export(outputFile, fileType)
	return err
}

func init() {
	exportCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Specify the output path of the snippet")
	exportCmd.Flags().StringVarP(&fileType, "type", "t", snippet.EXPORT_TYPE_CORGI, fmt.Sprintf("Choose export file type. Allowed values are: \"%s\", \"%s\".", snippet.EXPORT_TYPE_CORGI, snippet.EXPORT_TYPE_SHELL))
	rootCmd.AddCommand(exportCmd)
}
