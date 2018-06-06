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

type TemplateFieldMap map[string]*TemplateField // map from field name to template field object

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
	if err := snippet.writeToFile(filePath); err != nil {
		color.Red("Failure")
		return err
	}
	color.Green("Success")
	return nil
}

func (snippet *Snippet) Export(outputPath string) error {
	fmt.Printf("Exporting snippet %s... ", snippet.Title)
	if err := snippet.writeToFile(outputPath); err != nil {
		color.Red("Failure")
		return err
	}
	color.Green("Success")
	return nil
}

func (snippet *Snippet) writeToFile(filePath string) error {
	data, err := json.MarshalIndent(snippet, util.JSON_MARSHAL_PREFIX, util.JSON_MARSHAL_INDENT)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(filePath, data, 0644); err != nil {
		return err
	}
	return nil
}

func (snippet *Snippet) Execute() error {
	fmt.Println(color.GreenString("Start executing snippet \"%s\"...\n", snippet.Title))
	templateFieldMap := &TemplateFieldMap{}
	for idx, step := range snippet.Steps {
		stepCount := idx + 1
		fmt.Printf("%s: %s\n", color.GreenString("Step %d", stepCount), color.YellowString(step.Description))
		if err := step.Execute(templateFieldMap); err != nil {
			color.Red("[ Failure ]")
			return err
		}
		color.Green("[ Success ]")
		fmt.Println("")
	}
	return nil
}

func (snippet *Snippet) Describe() {
	fmt.Printf("%s: %s\n", color.YellowString("Title"), snippet.Title)
	for idx, step := range snippet.Steps {
		fmt.Printf("%s %s\n", color.YellowString("Step %d -", idx+1), step.Description)
	}
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

func (tfMap TemplateFieldMap) AddTemplateFieldIfNotExist(t *TemplateField) {
	if _, ok := tfMap[t.FieldName]; ok {
		// take the latest non-empty default value
		if t.Value != "" {
			tfMap[t.FieldName] = t
		}
	} else {
		tfMap[t.FieldName] = t

	}
}
