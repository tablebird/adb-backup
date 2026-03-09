package shell

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	whitespaceRegex = regexp.MustCompile(`^\s*$`)
)

func prepareCommandLine(cmd string, args ...string) (string, error) {
	if isBlank(cmd) {
		return "", errors.New("command cannot be empty")
	}

	for i, arg := range args {
		if strings.ContainsRune(arg, '"') {
			return "", fmt.Errorf("arg at index %d contains an invalid double quote: %s", i, arg)
		}
		if containsWhitespace(arg) {
			args[i] = fmt.Sprintf("\"%s\"", arg)
		}
	}

	// Prepend the command to the args array.
	if len(args) > 0 {
		cmd = fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
	}

	return cmd, nil
}

func containsWhitespace(str string) bool {
	return strings.ContainsAny(str, " \t\v")
}

func isBlank(str string) bool {
	return whitespaceRegex.MatchString(str)
}
