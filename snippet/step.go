package snippet

import (
	"corgi/util"
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type StepInfo struct {
	Command           string `json:"command"`
	Description       string `json:"description,omitempty"`
	ExecuteConcurrent bool   `json:"execute_concurrent"`
	templateFields    []*TemplateField
}

type TemplateField struct {
	FieldName    string `json:"field_name"`
	DefaultValue string `json:"default_value"`
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
	step.templateFields = ParseTemplateFields(cmd)
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
	fmt.Printf("%s: %s\n", color.GreenString("Running"), color.YellowString(step.Command))
	commandsList := strings.Split(step.Command, "&&")
	for _, c := range commandsList {
		c = strings.TrimSpace(c)
		cmdName := strings.Split(c, " ")[0]
		cmdArgs := strings.Split(c, " ")[1:]
		cmd := exec.Command(cmdName, cmdArgs...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func ParseTemplateFields(c string) []*TemplateField {
	re := regexp.MustCompile(`<([^(<>|\s)]+)>`)
	params := re.FindAllString(c, -1)
	templateFields := make([]*TemplateField, len(params), len(params))
	for idx, p := range params {
		// I'm doing this cuz I suck at building regex
		p = p[1 : len(p)-1]
		// fetch field and default value (if there's any)
		var field, defaultVal string
		if strings.Contains(p, "=") {
			field = strings.Split(p, "=")[0]
			defaultVal = strings.Split(p, "=")[1]
		} else {
			field = p
		}
		templateFields[idx] = &TemplateField{
			FieldName:    field,
			DefaultValue: defaultVal,
		}
	}
	return templateFields
}
