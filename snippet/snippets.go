package snippet

import (
	"corgi/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Snippets struct {
	Snippets []*jsonSnippet `json:"snippets"`
	fileLoc  string
}

type jsonSnippet struct {
	FileLoc string `json:"file_loc"`
	Title   string `json:"title"`
}

func LoadSnippets(filePath string) (*Snippets, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, err
	}
	snippets := &Snippets{}
	if err := util.LoadJsonDataFromFile(filePath, snippets); err != nil {
		return nil, err
	}
	if snippets.fileLoc == "" {
		snippets.fileLoc = filePath
	}
	return snippets, nil
}

func (snippets *Snippets) Save() error {
	if _, err := os.Stat(snippets.fileLoc); os.IsNotExist(err) {
		return err
	}
	data, err := json.Marshal(snippets)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(snippets.fileLoc, data, 0644); err != nil {
		return err
	}
	return nil
}

func (snippets *Snippets) AddSnippet(snippet *Snippet) {
	jsonSnippet := &jsonSnippet{
		Title:   snippet.Title,
		FileLoc: snippet.fileLoc,
	}
	snippets.Snippets = append(snippets.Snippets, jsonSnippet)
}

func (snippets *Snippets) FindSnippet(title string) (*Snippet, error) {
	for _, snp := range snippets.Snippets {
		if snp.Title == title {
			s, err := LoadSnippet(snp.FileLoc)
			if err != nil {
				return nil, err
			}
			return s, nil
		}
	}
	return nil, fmt.Errorf("could not find snippet with name: %s", title)
}
