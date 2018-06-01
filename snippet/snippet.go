package snippet

import (
	"fmt"
)

type Snippets struct {
	Snippets []Snippet `json:"snippets"`
}

type Snippet struct {
	Title   string     `json:"title"`
	Steps   []StepInfo `json:"steps"`
	FileLoc string     `json:"file_loc"`
}

type StepInfo struct {
	Command        string   `json:"command"`
	Description    string   `json:"description"`
	TemplateFields []string `json:"template_fields"`
}

func NewStepInfo(command string) *StepInfo {
	return &StepInfo{
		Command: command,
	}
}

func (step *StepInfo) WriteAnswer(field string, value interface{}) error {
	if field == "command" {
		step.Command = value.(string)
	} else if field == "description" {
		step.Description = value.(string)
	} else {
		return fmt.Errorf("unknown field %s", field)
	}
	return nil
}

func (step *StepInfo) AskQuestion(command string) error {
	return nil
}
