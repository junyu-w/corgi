package cmd

import (
	"bytes"
	"corgi/config"
	"corgi/snippet"
	"corgi/util"
	"errors"
	"io"
	"os"
	"strings"
	"github.com/fatih/color"
)

var MissingSnippetTitleError = errors.New("snippet title is not selected")

func loadConfigAndSnippetsMeta() (*config.Config, *snippet.SnippetsMeta, error) {
	// load config
	conf, err := config.Load()
	if err != nil {
		return nil, nil, err
	}
	// Load snippets
	snippets, err := snippet.LoadSnippetsMeta(conf.SnippetsFile)
	if err != nil {
		return nil, nil, err
	}
	return conf, snippets, nil
}

func filter(filterCmd string, candidates []string) (string, error) {
	var buf bytes.Buffer
	inputs := strings.Join(candidates, "\n")
	ws := io.MultiWriter(os.Stdout, &buf)
	if err := util.Execute(filterCmd, strings.NewReader(inputs), ws); err != nil {
		return "", err
	}
	result := strings.Trim(strings.TrimSpace(buf.String()), "\n")
	return result, nil
}

func filterSnippetTitle(filterCmd string, titles []string) (string, error) {
	if filterCmd != "" {
		title, err := filter(filterCmd, titles)
		if err != nil || title == "" {
			return "", MissingSnippetTitleError
		}
		return title, nil
	} else {
		color.Red("Install a fuzzy finder (\"fzf\" or \"peco\") to enable interactive selection")
		return "", MissingSnippetTitleError
	}
}
