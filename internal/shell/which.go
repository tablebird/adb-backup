package shell

import (
	"strings"
)

const _WHICH = "which"

func Which(s Shell, args ...string) (string, error) {
	return s.RunCommand(_WHICH, args...)
}

func WhichCommandExists(s Shell, command string) bool {
	res, err := Which(s, command)
	if err != nil {
		return false
	}
	return strings.TrimSpace(res) != ""
}
