package snippet

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/DrakeW/corgi/util"
	"github.com/fatih/color"
)

type StepInfo struct {
	Command     string `json:"command"`
	Description string `json:"description,omitempty"`
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

func (step *StepInfo) ConvertToShellScript(templates *TemplateFieldMap) string {
	shellCmds := make([]string, 0)
	templateFieldsMap := ParseTemplateFieldsMap(step.Command)
	for field := range templateFieldsMap {
		existingTf, _ := (*templates)[field]
		// add user input prompt
		if !existingTf.Asked {
			inputPromptShell := fmt.Sprintf("read -p \"%sEnter value for <%s>:%s\" %s", util.SHELL_GREEN, field, util.SHELL_NO_COLOR, field)
			existingTf.Asked = true
			shellCmds = append(shellCmds, inputPromptShell)
		}
	}
	// add command execution part
	re := regexp.MustCompile(TemplateParamsRegex)
	executeShell := re.ReplaceAllStringFunc(step.Command, func(sub string) string {
		field, _ := getParamNameAndValue(sub)
		return fmt.Sprintf("$%s", field)
	})
	runningPromptShell := fmt.Sprintf("echo -e \"%sRunning: %s%s%s\"", util.SHELL_GREEN, util.SHELL_YELLOW, executeShell, util.SHELL_NO_COLOR)
	shellCmds = append(shellCmds, runningPromptShell, executeShell)
	return strings.Join(shellCmds, "\n")
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
