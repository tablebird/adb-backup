package display

import (
	"adb-backup/internal/shell"
)

type DisplayManager interface {
	Size() (int, int)

	IsOn() bool
}

func NewDisplayManager(s shell.Shell) DisplayManager {
	return &shellDisplay{s: s}
}

type shellDisplay struct {
	s shell.Shell
}

func (d *shellDisplay) Size() (int, int) {
	width, height, _ := shell.WmSize(d.s)
	return width, height
}

func (d *shellDisplay) IsOn() bool {
	displayState, err := shell.DumpDisplayState(d.s)
	return err == nil && displayState
}
