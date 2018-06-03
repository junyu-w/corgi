package snippet

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"os"
	"strings"
)

const (
	SHELL_BASH        = "bash"
	SHELL_ZSH         = "zsh"
	SHELL_UNSUPPORTED = "unsupported"
)

var shellType string
var TempHistFile = "/tmp/corgi.hist"

func getHistoryFilePath() (string, error) {
	histFilePath, suc := os.LookupEnv("HISTFILE")
	if !suc {
		color.Red("Could not find HISTFILE in env - Using default based on shell type")
	}
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "zsh") {
		shellType = SHELL_ZSH
		histFilePath = fmt.Sprintf("%s/.zsh_history", homeDir)
	} else if strings.Contains(shell, "bash") {
		shellType = SHELL_BASH
		histFilePath = fmt.Sprintf("%s/.bash_history", homeDir)
	} else {
		shellType = SHELL_UNSUPPORTED
		return "", errors.New("only Bash, Zsh are currently supported.")
	}
	if _, err := os.Stat(histFilePath); os.IsNotExist(err) {
		return "", err
	}
	return histFilePath, nil
}

func ParseFileToStringArray(filePath string, parser CommandParser) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parsedLine := parser.Parse(line)
		lines = append(lines, parsedLine)
	}
	return lines, nil
}

func ReadShellHistory() ([]string, error) {
	histFilePath, err := getHistoryFilePath()
	if err != nil {
		return nil, err
	}
	parser, err := GetCmdParser(shellType)
	if err != nil {
		return nil, err
	}
	lines, err := ParseFileToStringArray(histFilePath, parser)
	return lines, nil
}

func SetUpHistFile(histCmds []string) error {
	// write commands to temp history file
	f, err := os.OpenFile(TempHistFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer f.Close()
	for _, cmd := range histCmds {
		if _, err := f.WriteString(fmt.Sprintf("%s\n", cmd)); err != nil {
			return err
		}
	}
	return nil
}

func RemoveHistFile() error {
	if err := os.Remove(TempHistFile); err != nil {
		return err
	}
	return nil
}
