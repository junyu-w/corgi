package snippet

import (
	"corgi/util"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
)

type SnippetsMeta struct {
	Snippets []*jsonSnippet `json:"snippets"`
	fileLoc  string
}

type jsonSnippet struct {
	FileLoc string `json:"file_loc"`
	Title   string `json:"title"`
}

func LoadSnippetsMeta(filePath string) (*SnippetsMeta, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, err
	}
	snippets := &SnippetsMeta{}
	if err := util.LoadJsonDataFromFile(filePath, snippets); err != nil {
		return nil, err
	}
	if snippets.fileLoc == "" {
		snippets.fileLoc = filePath
	}
	return snippets, nil
}

func (sm *SnippetsMeta) Save() error {
	if _, err := os.Stat(sm.fileLoc); os.IsNotExist(err) {
		return err
	}
	data, err := json.MarshalIndent(sm, util.JSON_MARSHAL_PREFIX, util.JSON_MARSHAL_INDENT)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(sm.fileLoc, data, 0644); err != nil {
		return err
	}
	return nil
}

// Save new snippet into snippetsDir and update snippets meta file
func (sm *SnippetsMeta) SaveNewSnippet(snippet *Snippet, snippetsDir string) error {
	if err := snippet.Save(snippetsDir); err != nil {
		return err
	}
	jsonSnippet := &jsonSnippet{
		Title:   snippet.Title,
		FileLoc: snippet.fileLoc,
	}
	sm.Snippets = append(sm.Snippets, jsonSnippet)
	if err := sm.Save(); err != nil {
		return err
	}
	return nil
}

func (sm *SnippetsMeta) DeleteSnippet(title string) error {
	idx, err := sm.findJsonSnippetIndex(title)
	if err != nil {
		return err
	}
	s := sm.Snippets[idx]
	fmt.Printf("Deleting snippet %s... ", s.Title)
	// delete snippet file
	fileLoc := s.FileLoc
	if err = os.Remove(fileLoc); err != nil {
		color.Red("Failure")
		return err
	}
	// delete from snippets meta
	sm.Snippets = append(sm.Snippets[:idx], sm.Snippets[idx+1:]...)
	if err = sm.Save(); err != nil {
		color.Red("Failure")
		return err
	}
	color.Green("Success")
	return nil
}

func (sm *SnippetsMeta) FindSnippet(title string) (*Snippet, error) {
	idx, err := sm.findJsonSnippetIndex(title)
	if err != nil {
		return nil, err
	}
	s, err := LoadSnippet(sm.Snippets[idx].FileLoc)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (sm *SnippetsMeta) findJsonSnippetIndex(title string) (int, error) {
	idx := -1
	for i, snp := range sm.Snippets {
		if snp.Title == title {
			idx = i
			break
		}
	}
	if idx == -1 {
		return idx, fmt.Errorf("could not find snippet with name: %s", title)
	}
	return idx, nil
}

func (sm *SnippetsMeta) GetSnippetTitles() []string {
	titles := make([]string, len(sm.Snippets))
	for idx, s := range sm.Snippets {
		titles[idx] = s.Title
	}
	return titles
}
