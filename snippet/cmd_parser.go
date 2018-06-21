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
type FishCmdParser struct{}

func (z ZshCmdParser) Parse(line string) string {
	parts := strings.Split(line, ";")
	return strings.Join(parts[1:], ";")
}

func (b BashCmdParser) Parse(line string) string {
	return line
}

func (f FishCmdParser) Parse(line string) string {
	fishCmdPrefix := "- cmd: "
	if strings.HasPrefix(line, fishCmdPrefix) {
		return line[len(fishCmdPrefix):]
	}
	return ""
}

func GetCmdParser(shellType string) (CommandParser, error) {
	if shellType == SHELL_ZSH {
		return ZshCmdParser{}, nil
	} else if shellType == SHELL_BASH {
		return BashCmdParser{}, nil
	} else if shellType == SHELL_FISH {
		return FishCmdParser{}, nil
	}
	return nil, fmt.Errorf("unsupported shell type: %s", shellType)
}
