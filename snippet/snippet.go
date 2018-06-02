package snippet

import (
	"encoding/json"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/kataras/iris/core/errors"
	"io/ioutil"
	"strings"
)

type Snippet struct {
	Title   string      `json:"title"`
	Steps   []*StepInfo `json:"steps"`
	FileLoc string      `json:"file_loc"`
}

type StepInfo struct {
	Command           string   `json:"command"`
	Description       string   `json:"description,omitempty"`
	executeConcurrent bool     `json:"execute_concurrent"`
	TemplateFields    []string `json:"template_fields"`
}

type Answerable interface {
	AskQuestion(options ...interface{}) error
}

func scan(prompt string, defaultInp string) (string, error) {
	// create config
	config := &readline.Config{
		Prompt:            prompt,
		HistoryFile:       TempHistFile,
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

// ################### Step related code ############################

func NewStepInfo(command string) *StepInfo {
	return &StepInfo{
		Command: command,
	}
}

func (step *StepInfo) AskQuestion(options ...interface{}) error {
	// set command
	cmd, err := scan(color.GreenString("Command: "), step.Command)
	if err != nil {
		return err
	}
	// TODO: read template from command
	step.Command = cmd
	// set description
	description, err := scan(color.GreenString("Description: "), "")
	if err != nil {
		return err
	}
	step.Description = description
	return nil
}

// ################### Snippet related code ############################

func NewSnippet(title string, cmds []string) (*Snippet, error) {
	snippet := &Snippet{
		Title: title,
	}
	if err := snippet.AskQuestion(cmds); err != nil {
		return nil, err
	}
	return snippet, nil
}

func (snippet *Snippet) AskQuestion(options ...interface{}) error {
	// ask about each step
	initialDefaultCmds := options[0].([]string)
	stepCount := 0
	steps := make([]*StepInfo, 0)
	for {
		color.Yellow("Step %d:", stepCount+1)
		var defaultCmd string
		if stepCount < len(initialDefaultCmds) {
			defaultCmd = initialDefaultCmds[stepCount]
		}
		step := NewStepInfo(defaultCmd)
		err := step.AskQuestion()
		if err != nil {
			return err
		}
		steps = append(steps, step)
		var addOneMoreStep bool
		for {
			addStepInp, err := scan(color.RedString("Add another step? (y/n): "), "")
			if err != nil {
				return err
			}
			if addStepInp == "y" {
				addOneMoreStep = true
			} else if addStepInp == "n" {
				addOneMoreStep = false
			} else {
				continue
			}
			break
		}
		fmt.Print("\n")
		if !addOneMoreStep {
			break
		}
		stepCount++
	}
	snippet.Steps = steps
	// ask about title if not set
	if snippet.Title == "" {
		title, err := scan(color.YellowString("Title: "), "")
		if err != nil {
			return err
		}
		snippet.Title = title
	}
	return nil
}

func (snippet *Snippet) Save(snippetsDir string) error {
	fmt.Printf("Saving snippet %s... ", snippet.Title)
	data, err := json.Marshal(snippet)
	if err != nil {
		color.RedString("Failed")
		return err
	}
	filePath := fmt.Sprintf("%s/%s.json", snippetsDir, snippet.Title)
	if err = ioutil.WriteFile(filePath, data, 0644); err != nil {
		color.RedString("Failed")
		return err
	}
	color.GreenString("Success")
	return nil
}
