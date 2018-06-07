package util

import (
	"encoding/json"
	"errors"
	"github.com/chzyer/readline"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const (
	JSON_MARSHAL_PREFIX = ""
	JSON_MARSHAL_INDENT = "  "
	STEP_RANGE_SEP      = "-"
)

func LoadJsonDataFromFile(filePath string, object interface{}) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, object); len(data) > 0 && err != nil {
		return err
	}
	return nil
}

func Scan(prompt string, defaultInp string, historyFile string) (string, error) {
	// create config
	config := &readline.Config{
		Prompt:            prompt,
		HistoryFile:       historyFile,
		HistorySearchFold: true,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
	}
	rl, err := readline.NewEx(config)
	if err != nil {
		return "", err
	}
	defer rl.Close()

	for {
		line, err := rl.ReadlineWithDefault(defaultInp)
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		return line, nil
	}
	return "", errors.New("cancelled")
}

func Execute(command string, r io.Reader, w io.Writer) error {
	cmd := exec.Command("sh", "-c", strings.TrimSpace(command))
	cmd.Stdin = r
	cmd.Stdout = w
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
