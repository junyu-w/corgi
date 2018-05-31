package snippet

type Snippets struct {
	Snippets []Snippet
}

type Snippet struct {
	Title   string
	Steps   []StepInfo
	FileLoc string
}

type StepInfo struct {
	Command        string
	Description    string
	TemplateFields []string
}
