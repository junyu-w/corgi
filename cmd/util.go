package cmd

import (
	"bytes"
	"corgi/config"
	"corgi/snippet"
	"corgi/util"
	"io"
	"os"
	"strings"
)

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
