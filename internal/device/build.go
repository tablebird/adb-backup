package device

import (
	"adb-backup/internal/shell"
)

type Build interface {
	VersionRelease() int
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
