package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	SnippetsFile string `json:"snippets_file"`
	SnippetsDir  string `json:"snippets_dir"`
}

const (
	DEFAULT_CONFIG_FILE   = ".corgi/corgi_conf.json"
	DEFAULT_SNIPPETS_DIR  = ".corgi/snippets"
	DEFAULT_SNIPPETS_FILE = ".corgi/snippets.json"
)

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

func Load() (*Config, error) {
	configFile, err := GetDefaultConfigFile()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	if err = json.Unmarshal(data, config); len(data) > 0 && err != nil {
		return nil, err
	}
	// if config file has not content, initialize it with default
	if config.SnippetsFile == "" && config.SnippetsDir == "" {
		snippetsFile, err := GetDefaultSnippetsFile()
		if err != nil {
			return nil, err
		}
		config.SnippetsFile = snippetsFile
		snippetsDir, err := GetDefaultSnippetsDir()
		if err != nil {
			return nil, err
		}
		config.SnippetsDir = snippetsDir
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
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(confFile, data, 0644); err != nil {
		return err
	}
	return nil
}
