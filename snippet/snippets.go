package snippet

import (
	"corgi/util"
	"encoding/json"
	"io/ioutil"
	"os"
)

type Snippets struct {
	Snippets []*Snippet `json:"snippets"`
	fileLoc  string     `json:"file_loc"`
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
