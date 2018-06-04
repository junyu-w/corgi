package snippet

import (
	"corgi/util"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type StepInfo struct {
	Command           string `json:"command"`
	Description       string `json:"description,omitempty"`
	ExecuteConcurrent bool   `json:"execute_concurrent"`
}

var TemplateParamsRegex = `<([^(<>|\s)]+)>`

type TemplateField struct {
	FieldName string `json:"field_name"`
	Value     string `json:"default_value"`
}

type TemplateFieldMap map[string]*TemplateField

func NewStepInfo(command string) *StepInfo {
	return &StepInfo{
		Command: command,
	}
}

func (step *StepInfo) AskQuestion(options ...interface{}) error {
	// set command
	cmd, err := util.Scan(color.GreenString("Command: "), step.Command, TempHistFile)
	if err != nil {
		return err
	}
	step.Command = cmd
	// set description
	description, err := util.Scan(color.GreenString("Description: "), "", TempHistFile)
	if err != nil {
		return err
	}
	step.Description = description
	return nil
}

// TODO: add concurrent execution
func (step *StepInfo) Execute() error {
	// fill in templates
	templateFieldsMap := ParseTemplateFields(step.Command)
	for _, t := range templateFieldsMap {
		t.AskQuestion()
	}
	// replace params in command with input values
	command := FillTemplates(step.Command, &templateFieldsMap)
	// execute command
	fmt.Printf("%s: %s\n", color.GreenString("Running"), color.YellowString(command))
	cmd := exec.Command("sh", "-c", strings.TrimSpace(command))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// getParamNameAndValue returns name and value of param, parma must be
// a string that matches TemplateParamsRegex
func getParamNameAndValue(p string) (string, string) {
	// I'm doing this cuz I suck at building regex
	p = p[1 : len(p)-1]
	// fetch field and default value (if there's any)
	var field, val string
	if strings.Contains(p, "=") {
		field = strings.Split(p, "=")[0]
		val = strings.Split(p, "=")[1]
	} else {
		field = p
	}
	return field, val
}

func ParseTemplateFields(c string) TemplateFieldMap {
	re := regexp.MustCompile(TemplateParamsRegex)
	params := re.FindAllString(c, -1)
	tfMap := TemplateFieldMap{}
	for _, p := range params {
		field, defaultVal := getParamNameAndValue(p)
		tfMap.AddTemplateField(&TemplateField{
			FieldName: field,
			Value:     defaultVal,
		})
	}
	return tfMap
}

func FillTemplates(c string, tfMap *TemplateFieldMap) string {
	re := regexp.MustCompile(TemplateParamsRegex)
	filledCmd := re.ReplaceAllStringFunc(c, func(sub string) string {
		field, _ := getParamNameAndValue(sub)
		if t, ok := (*tfMap)[field]; ok {
			return t.Value
		}
		log.Panic(color.RedString("Couldn't find field with name %s", field))
		return ""
	})
	return filledCmd
}

func (tf *TemplateField) AskQuestion(options ...interface{}) error {
	val, err := util.Scan(color.CyanString("Enter value for <%s>: ", tf.FieldName), tf.Value, "")
	if err != nil {
		return err
	}
	tf.Value = val
	return nil
}

func (tfMap TemplateFieldMap) AddTemplateField(t *TemplateField) {
	if _, ok := tfMap[t.FieldName]; ok {
		// take the latest non-empty default value
		if t.Value != "" {
			tfMap[t.FieldName] = t
		}
	} else {
		tfMap[t.FieldName] = t
	}
}
