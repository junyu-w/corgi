package cmd

import (
	"corgi/config"
	"corgi/snippet"
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
