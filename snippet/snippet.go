package snippet

import (
	"fmt"
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

func (step *StepInfo) askQuestion() error {
	// TODO: finish this
	return nil
}

func (snippet *Snippet) AskQuestions(commands []string) error {
	// TODO: finish this
	fmt.Println("Asking Questions")
	// set up history file and readline instance

	// ask about each step
	for _, step := range snippet.Steps {
		err := step.askQuestion()
		if err != nil {
			return err
		}
	}
	// ask about title
	return nil
}

func (Snippet *Snippet) Save() error {
	// TODO: finish this
	fmt.Println("Saving snippet")
	return nil
}

func NewSnippet(commands []string) (*Snippet, error) {
	snippet := &Snippet{}
	if err := snippet.AskQuestions(commands); err != nil {
		return nil, err
	}
	return snippet, nil
}
