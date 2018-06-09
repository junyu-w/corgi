package snippet

import (
	"corgi/util"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"regexp"
	"strings"
)

type StepInfo struct {
	Command           string `json:"command"`
	Description       string `json:"description,omitempty"`
}

var TemplateParamsRegex = `<([^(<>|\s)]+)>`

type TemplateField struct {
	FieldName string
	Value     string
	Asked     bool
}

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
// valid options include 'useDefaultVal' indicated by the --use-default flag
func (step *StepInfo) Execute(templates *TemplateFieldMap, options ...interface{}) error {
	useDefaultVal := options[0].(bool)
	if !useDefaultVal {
		// fill in templates
		templateFieldsMap := ParseTemplateFieldsMap(step.Command)
		for field := range templateFieldsMap {
			existingTf, _ := (*templates)[field]
			// only ask once for user input for the same template field
			if !existingTf.Asked {
				existingTf.AskQuestion()
			}
		}
	}
	// replace params in command with input values
	command := FillTemplates(step.Command, templates)
	// execute command
	fmt.Printf("%s: %s\n", color.GreenString("Running"), color.YellowString(command))
	if err := util.Execute(command, os.Stdin, os.Stdout); err != nil {
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

func ParseTemplateFieldsMap(c string) TemplateFieldMap {
	re := regexp.MustCompile(TemplateParamsRegex)
	params := re.FindAllString(c, -1)
	tfMap := TemplateFieldMap{}
	for _, p := range params {
		field, defaultVal := getParamNameAndValue(p)
		tfMap.AddTemplateFieldIfNotExist(&TemplateField{
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
	val, err := util.Scan(color.GreenString("Enter value for <%s>: ", tf.FieldName), tf.Value, "")
	if err != nil {
		return err
	}
	tf.Value = val
	tf.Asked = true
	return nil
}
