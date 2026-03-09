package input

import (
	"adb-backup/internal/log"
	"adb-backup/internal/shell"
)

type InputManager interface {
	KeyEvent(keyCode ...int)
	Text(text string)
	Power()
}

func NewInputManager(s shell.Shell) InputManager {
	return &shellInput{s: s}
}

type shellInput struct {
	s shell.Shell
}

func (i *shellInput) Power() {
	err := shell.InputKeyEventPower(i.s)
	if err != nil {
		log.WarningF("InputKeyEventPower error %s", err.Error())
	}
}

func (i *shellInput) KeyEvent(keyCode ...int) {
	var err error
	if len(keyCode) == 1 {
		err = shell.InputKeyEvent(i.s, keyCode[0])
	} else {
		err = shell.InputKeyCombination(i.s, keyCode...)
	}
	if err != nil {
		log.WarningF("InputKeyEvent error %s", err.Error())
	}
}

func (i *shellInput) Text(text string) {
	err := shell.InputText(i.s, text)
	if err != nil {
		log.WarningF("InputText error %s", err.Error())
	}
}
