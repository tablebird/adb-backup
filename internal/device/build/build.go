package build

import (
	"adb-backup/internal/shell"
)

type Build interface {
	VersionRelease() int
}

func NewBuild(s shell.Shell) Build {
	return &shellBuild{s: s}
}

type shellBuild struct {
	s shell.Shell
}

func (b *shellBuild) VersionRelease() int {
	result, err := shell.GetPropBuildVersionRelease(b.s)
	if err == nil {
		return 0
	}
	return result
}
