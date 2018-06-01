package snippet

import (
	"fmt"
	"strings"
)

type CommandParser interface {
	Parse(string) string
}

type BashCmdParser struct{}
type ZshCmdParser struct{}

func (z ZshCmdParser) Parse(line string) string {
	parts := strings.Split(line, ";")
	return strings.Join(parts[1:], ";")
}

func (b BashCmdParser) Parse(line string) string {
	return line
}

func GetCmdParser(shellType string) (CommandParser, error) {
	if shellType == SHELL_ZSH {
		return ZshCmdParser{}, nil
	} else if shellType == SHELL_BASH {
		return BashCmdParser{}, nil
	} else {
		return nil, fmt.Errorf("unsupported shell type: %s", shellType)
	}
}
