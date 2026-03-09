package shell

import (
	"fmt"
	"strconv"
	"strings"
)

const _WM = "wm"

func Wm(s Shell, args ...string) (string, error) {
	return s.RunCommand(_WM, args...)
}

func WmSize(s Shell) (int, int, error) {
	res, err := Wm(s, "size")
	if err != nil {
		return 0, 0, err
	}
	split := strings.SplitN(res, ":", 2)
	if len(split) < 2 {
		return 0, 0, fmt.Errorf("wm size error")
	}
	split = strings.SplitN(strings.TrimSpace(split[1]), "x", 2)
	if len(split) < 2 {
		return 0, 0, fmt.Errorf("wm size error")
	}
	width, err := strconv.Atoi(strings.TrimSpace(split[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("width %s", err.Error())
	}
	height, err := strconv.Atoi(strings.TrimSpace(split[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("height %s", err.Error())
	}
	return width, height, nil
}
