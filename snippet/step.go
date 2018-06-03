package snippet

import (
	"corgi/util"
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"strings"
)

type StepInfo struct {
	Command           string   `json:"command"`
	Description       string   `json:"description,omitempty"`
	ExecuteConcurrent bool     `json:"execute_concurrent"`
	TemplateFields    []string `json:"template_fields"`
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
	// TODO: read template from command
	step.Command = cmd
	// set description
	description, err := util.Scan(color.GreenString("Description: "), "", TempHistFile)
	if err != nil {
		return err
	}
	step.Description = description
	return nil
}

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
