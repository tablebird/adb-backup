package shell

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	_INPUT = "input"

	_KEYEVENT = "keyevent"

	_TEXT = "text"

	_SWIPE = "swipe"

	_KEYCOMBINATION = "keycombination"
)

func Input(s Shell, args ...string) error {
	res, err := s.RunCommand(_INPUT, args...)
	if err != nil {
		return err
	}
	if strings.TrimSpace(res) != "" {
		return fmt.Errorf("input error: %s", res)
	}
	return nil
}

func InputCommand(s Shell, cmd string, args ...int) error {
	stringArgs := make([]string, len(args)+1)
	stringArgs[0] = cmd
	for i, arg := range args {
		stringArgs[i+1] = strconv.Itoa(arg)
	}
	return Input(s, stringArgs...)
}

func InputKeyEvent(s Shell, keyCode int) error {
	return InputCommand(s, _KEYEVENT, keyCode)
}

func InputKeyCombination(s Shell, keycode ...int) error {
	return InputCommand(s, _KEYCOMBINATION, keycode...)
}

func InputText(s Shell, text string) error {
	return Input(s, _TEXT, text)
}

func InputSwipe(s Shell, x1, y1, x2, y2 int, duration int) error {
	return InputCommand(s, _SWIPE, x1, y1, x2, y2, duration)
}

func InputKeyEventPower(s Shell) error {
	return InputKeyEvent(s, 26)
}
