package config

import (
	"corgi/util"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

type Config struct {
	SnippetsFile string `json:"snippets_file"`
	SnippetsDir  string `json:"snippets_dir"`
	Editor       string `json:"editor"`
	FilterCmd    string `json:"filter_cmd"`
}

const (
	DEFAULT_CONFIG_FILE     = ".corgi/corgi_conf.json"
	DEFAULT_SNIPPETS_DIR    = ".corgi/snippets"
	DEFAULT_SNIPPETS_FILE   = ".corgi/snippets.json"
	DEFAULT_EDITOR          = "vim"
	DEFAULT_FILTER_CMD_FZF  = "fzf"
	DEFAULT_FILTER_CMD_PECO = "peco"
)

var MissingDefaultFilterCmdError = errors.New("missing default filter cmd")

func getOrCreatePath(loc string, perm os.FileMode, isDir bool) error {
	dirPath := path.Dir(loc)
	if isDir {
		dirPath = loc
	}
	if _, err := os.Stat(loc); os.IsNotExist(err) {
		if err = os.MkdirAll(dirPath, perm); err != nil {
			return err
		}
		if !isDir {
			f, err := os.Create(loc)
			if err != nil {
				return err
			}
			defer f.Close()
		}
	}
	return nil
}

func GetDefaultConfigFile() (string, error) {
	var defaultConfigFileLoc = fmt.Sprintf("%s/%s", os.Getenv("HOME"), DEFAULT_CONFIG_FILE)
	if err := getOrCreatePath(defaultConfigFileLoc, 0755, false); err != nil {
		return "", err
	}
	return defaultConfigFileLoc, nil
}

func GetDefaultSnippetsDir() (string, error) {
	var defaultSnippetsDir = fmt.Sprintf("%s/%s", os.Getenv("HOME"), DEFAULT_SNIPPETS_DIR)
	if err := getOrCreatePath(defaultSnippetsDir, 0755, true); err != nil {
		return "", err
	}
	return defaultSnippetsDir, nil
}

func GetDefaultSnippetsFile() (string, error) {
	var defaultSnippetsFile = fmt.Sprintf("%s/%s", os.Getenv("HOME"), DEFAULT_SNIPPETS_FILE)
	if err := getOrCreatePath(defaultSnippetsFile, 0755, false); err != nil {
		return "", err
	}
	return defaultSnippetsFile, nil
}

func GetDefaultEditor() (string, error) {
	editorPath, suc := os.LookupEnv("EDITOR")
	if !suc {
		editorPath, err := exec.LookPath(DEFAULT_EDITOR)
		if err != nil {
			return "", fmt.Errorf("could not find %s (default) in $PATH, update your editor choice with \"corgi configure --editor <path to your editor>\"", DEFAULT_EDITOR)
		}
		return editorPath, nil
	}
	return editorPath, nil
}

func GetDefaultFilterCmd() (string, error) {
	filterCmdPath, err := exec.LookPath(DEFAULT_FILTER_CMD_PECO)
	if err != nil {
		filterCmdPath = ""
	}
	filterCmdPath, err = exec.LookPath(DEFAULT_FILTER_CMD_FZF)
	if err != nil {
		filterCmdPath = ""
	}
	if filterCmdPath == "" {
		return "", MissingDefaultFilterCmdError
	}
	return filterCmdPath, nil
}

func Load() (*Config, error) {
	configFile, err := GetDefaultConfigFile()
	if err != nil {
		return nil, err
	}
	config := &Config{}
	if err = util.LoadJsonDataFromFile(configFile, config); err != nil {
		return nil, err
	}
	// if config file has no content, initialize it with default
	if config.IsNew() {
		// set default snippets file
		snippetsFile, err := GetDefaultSnippetsFile()
		if err != nil {
			return nil, err
		}
		config.SnippetsFile = snippetsFile
		// set default snippets dir
		snippetsDir, err := GetDefaultSnippetsDir()
		if err != nil {
			return nil, err
		}
		config.SnippetsDir = snippetsDir
		// set default editor
		editor, err := GetDefaultEditor()
		if err != nil {
			return nil, err
		}
		config.Editor = editor
		// set default filter cmd
		filterCmd, err := GetDefaultFilterCmd()
		if err != nil && err != MissingDefaultFilterCmdError {
			return nil, err
		}
		config.FilterCmd = filterCmd
		// save
		config.Save()
	}
	return config, nil
}

func (c *Config) Save() error {
	// get config file
	confFile, err := GetDefaultConfigFile()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, util.JSON_MARSHAL_PREFIX, util.JSON_MARSHAL_INDENT)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(confFile, data, 0644); err != nil {
		return err
	}
	return nil
}

func (c *Config) IsNew() bool {
	return c.SnippetsFile == "" && c.SnippetsDir == "" && c.Editor == "" && c.FilterCmd == ""
}
