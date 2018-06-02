package snippet

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/kataras/iris/core/errors"
	"os"
)

type Snippets struct {
	Snippets []Snippet `json:"snippets"`
}

type Snippet struct {
	Title   string      `json:"title"`
	Steps   []*StepInfo `json:"steps"`
	FileLoc string      `json:"file_loc"`
}

type StepInfo struct {
	Command        string   `json:"command"`
	Description    string   `json:"description,omitempty"`
	TemplateFields []string `json:"template_fields"`
}

var tempHistFile = "/tmp/corgi.hist"

func setUpHistFile(histCmds []string) error {
	// write commands to temp history file
	f, err := os.OpenFile(tempHistFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, cmd := range histCmds {
		if _, err := f.WriteString(cmd); err != nil {
			return err
		}
	}
	return nil
}

func scan(prompt string, defaultInp string) (string, error) {
	// create config
	config := &readline.Config{
		Prompt:            prompt,
		HistoryFile:       tempHistFile,
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

func (step *StepInfo) askQuestion() error {
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

func NewSnippet(commands []string) (*Snippet, error) {
	snippet := &Snippet{}
	if err := snippet.AskQuestions(commands); err != nil {
		return nil, err
	}
	return snippet, nil
}

func (snippet *Snippet) AskQuestions(commands []string) error {
	// TODO: use full command history
	// set up history file
	if err := setUpHistFile(commands); err != nil {
		return err
	}
	// ask about each step
	steps := make([]*StepInfo, len(commands))
	// TODO: make step adding process interactive - asking if user want to add a step
	for idx, cmd := range commands {
		color.Yellow("Step %d:", idx+1)
		step := NewStepInfo(cmd)
		err := step.askQuestion()
		if err != nil {
			return err
		}
		steps = append(steps, step)
		fmt.Println("")
	}
	snippet.Steps = steps
	// ask about title
	title, err := scan(color.YellowString("Title: "), "")
	if err != nil {
		return err
	}
	snippet.Title = title
	return nil
}

func (Snippet *Snippet) Save() error {
	// TODO: finish this
	fmt.Println("Saving snippet")
	return nil
}
