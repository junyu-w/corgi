package snippet

import (
	"corgi/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"strconv"
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

var (
	MissingDefaultValueError = errors.New("missing default value for template field")
	InvalidStepRangeError    = errors.New("step range specified is invalid")
)

func NewSnippet(title string, cmds []string) (*Snippet, error) {
	snippet := &Snippet{
		Title: title,
	}
	if err := snippet.AskQuestion(cmds); err != nil {
		return nil, err
	}
	return snippet, nil
}

func LoadSnippet(filePath string) (*Snippet, error) {
	snippet := &Snippet{}
	if err := util.LoadJsonDataFromFile(filePath, snippet); err != nil {
		return nil, err
	}
	snippet.fileLoc = filePath
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

func getSnippetFileName(title string) string {
	return fmt.Sprintf("%s.json", strings.Replace(title, " ", "_", -1))
}

func (snippet *Snippet) Save(snippetsDir string) error {
	fmt.Printf("Saving snippet \"%s\"... ", snippet.Title)
	filePath := fmt.Sprintf("%s/%s", snippetsDir, getSnippetFileName(snippet.Title))
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

func (snippet *Snippet) Execute(options ...interface{}) error {
	fmt.Println(color.GreenString("Start executing snippet \"%s\"...\n", snippet.Title))
	// build template fields
	useDefaultVal := options[0].(bool)
	templateFieldMap := snippet.BuildTemplateFieldMap()
	if useDefaultVal {
		// check if all fields has default if --use-default set
		fieldWithNoDefault := make([]string, 0)
		for field, tf := range templateFieldMap {
			if tf.Value == "" {
				fieldWithNoDefault = append(fieldWithNoDefault, fmt.Sprintf("<%s>", field))
			}
		}
		if len(fieldWithNoDefault) > 0 {
			color.Red("[ Failure ] - Template field(s) %s do(es) not have default value set", strings.Join(fieldWithNoDefault, ", "))
			return MissingDefaultValueError
		}
	}
	// select step range
	stepRange := options[1].(string)
	start, end, err := snippet.ParseStepRangeToIdx(stepRange)
	if err != nil {
		color.Red("[ Failure ] - %s", err.Error())
		return err
	}
	for idx, step := range snippet.Steps[start:end] {
		stepCount := start + idx + 1
		fmt.Printf("%s: %s\n", color.GreenString("Step %d", stepCount), color.YellowString(step.Description))
		if err := step.Execute(&templateFieldMap, useDefaultVal); err != nil {
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

func (snippet *Snippet) BuildTemplateFieldMap() TemplateFieldMap {
	tfMap := TemplateFieldMap{}
	for _, step := range snippet.Steps {
		curTfMap := ParseTemplateFieldsMap(step.Command)
		for _, tf := range curTfMap {
			tfMap.AddTemplateFieldIfNotExist(tf)
		}
	}
	return tfMap
}

func (snippet *Snippet) ParseStepRangeToIdx(stepRange string) (int, int, error) {
	if stepRange == "" {
		return 0, len(snippet.Steps), nil
	}
	if strings.Contains(stepRange, util.STEP_RANGE_SEP) {
		sRange := strings.Split(stepRange, util.STEP_RANGE_SEP)
		if sRange[0] == "" {
			return -1, -1, InvalidStepRangeError
		} else if sRange[1] == "" {
			start, err := strconv.ParseInt(sRange[0], 10, 32)
			if err != nil {
				return -1, -1, err
			}
			end := len(snippet.Steps)
			startIdx := int(start) - 1
			endIdx := end
			// check for validity
			if isStepRangeInvalid(startIdx, endIdx, len(snippet.Steps)) {
				return -1, -1, InvalidStepRangeError
			}
			return startIdx, endIdx, nil
		} else {
			start, err := strconv.ParseInt(sRange[0], 10, 32)
			if err != nil {
				return -1, -1, err
			}
			end, err := strconv.ParseInt(sRange[1], 10, 32)
			if err != nil {
				return -1, -1, err
			}
			startIdx := int(start) - 1
			endIdx := int(end)
			if isStepRangeInvalid(startIdx, endIdx, len(snippet.Steps)) {
				return -1, -1, InvalidStepRangeError
			}
			return startIdx, endIdx, nil
		}
	} else {
		start, err := strconv.ParseInt(stepRange, 10, 32)
		if err != nil {
			return -1, -1, err
		}
		startIdx := int(start) - 1
		endIdx := int(start)
		if isStepRangeInvalid(startIdx, endIdx, len(snippet.Steps)) {
			return -1, -1, InvalidStepRangeError
		}
		return startIdx, endIdx, nil
	}
}

func isStepRangeInvalid(start, end, length int) bool {
	return start < 0 || end > length || start >= end
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
