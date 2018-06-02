package snippet

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/kataras/iris/core/errors"
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
	if _, err := os.Stat(histFilePath); err != nil {
		return "", err
	}
	return histFilePath, nil
}

func ReadShellHistory(last int) ([]string, error) {
	last += 1 // we want at least the command before the corgi command itself
	histFilePath, err := getHistoryFilePath()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(histFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	parser, err := GetCmdParser(shellType)
	if err != nil {
		return nil, err
	}
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parsedLine := parser.Parse(line)
		lines = append(lines, parsedLine)
	}
	return lines[len(lines)-last : len(lines)-1], nil
}
