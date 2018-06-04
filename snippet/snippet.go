package snippet

import (
	"corgi/util"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"strings"
)

type Snippet struct {
	Title   string      `json:"title"`
	Steps   []*StepInfo `json:"steps"`
	fileLoc string
}

type Answerable interface {
	AskQuestion(options ...interface{}) error
}

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
	// check options
	initialDefaultCmds := options[0].([]string)
	// ask about each step
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
			addStepInp, err := util.Scan(color.RedString("Add another step? (y/n): "), "", TempHistFile)
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
		title, err := util.Scan(color.YellowString("Title: "), "", TempHistFile)
		if err != nil {
			return err
		}
		snippet.Title = title
	}
	return nil
}

func (snippet *Snippet) Save(snippetsDir string) error {
	fmt.Printf("Saving snippet %s... ", snippet.Title)
	filePath := fmt.Sprintf("%s/%s.json", snippetsDir, strings.Replace(snippet.Title, " ", "_", -1))
	snippet.fileLoc = filePath
	data, err := json.Marshal(snippet)
	if err != nil {
		color.Red("Failure")
		return err
	}
	if err = ioutil.WriteFile(filePath, data, 0644); err != nil {
		color.Red("Failure")
		return err
	}
	color.Green("Success")
	return nil
}

func (snippet *Snippet) Execute() error {
	fmt.Println(color.GreenString("Start executing snippet \"%s\"...\n", snippet.Title))
	for idx, step := range snippet.Steps {
		stepCount := idx + 1
		fmt.Printf("%s: %s\n", color.GreenString("Step %d", stepCount), color.YellowString(step.Description))
		if err := step.Execute(); err != nil {
			color.Red("[ Failure ]")
			return err
		}
		color.Green("[ Success ]")
		fmt.Println("")
	}
	return nil
}

func (snippet *Snippet) GetFilePath() string {
	return snippet.fileLoc
}

func LoadSnippet(filePath string) (*Snippet, error) {
	snippet := &Snippet{}
	if err := util.LoadJsonDataFromFile(filePath, snippet); err != nil {
		return nil, err
	}
	snippet.fileLoc = filePath
	return snippet, nil
}
