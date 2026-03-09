package screen

import (
	"adb-backup/internal/shell"
	"io"
)

type ScreenManager interface {
	LiveH264() (io.ReadCloser, error)
	Capture() (io.ReadCloser, error)
}

func NewScreenManager(d shell.ReaderCloserShell) ScreenManager {
	return &shellScreenRecord{s: d}
}

type shellScreenRecord struct {
	s shell.ReaderCloserShell
}

func (s *shellScreenRecord) LiveH264() (io.ReadCloser, error) {
	return shell.ScreenRecordH264Live(s.s, "--bit-rate=4000000")
}

func (s *shellScreenRecord) Capture() (io.ReadCloser, error) {
	return shell.ScreenCap(s.s, "-p")
}
