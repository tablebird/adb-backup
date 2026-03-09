package device

import (
	"adb-backup/internal/shell"

	adb "github.com/zach-klippenstein/goadb"
)

type Build interface {
	VersionRelease() int
}

type shellBuild struct {
	adbDevice *adb.Device
}

func (b *shellBuild) VersionRelease() int {
	result, err := shell.GetPropBuildVersionRelease(b.adbDevice)
	if err == nil {
		return 0
	}
	return result
}
