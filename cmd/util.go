package cmd

import (
	"corgi/config"
	"corgi/snippet"
)

func loadConfigAndSnippets() (*config.Config, *snippet.Snippets, error) {
	// load config
	conf, err := config.Load()
	if err != nil {
		return nil, nil, err
	}
	// Load snippets
	snippets, err := snippet.LoadSnippets(conf.SnippetsFile)
	if err != nil {
		return nil, nil, err
	}
	return conf, snippets, nil
}
