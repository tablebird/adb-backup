package device

import (
	"adb-backup/internal/log"
	"adb-backup/internal/shell"
)

type PowerManager interface {
	BatteryLevel() int
	CharingType() []BatteryPoweredType
}

type shellPower struct {
	s shell.Shell
}

func (p *shellPower) BatteryLevel() int {
	level, err := shell.DumpBatteryLevel(p.s)
	if err != nil {
		return 0
	}
	return level
}

func (p *shellPower) CharingType() []BatteryPoweredType {
	charging, err := shell.DumpBatteryPoweredType(p.s)
	if err != nil {
		return nil
	}
	var result []BatteryPoweredType
	for _, c := range charging {
		t, e := parseBatteryPoweredType(c)
		if e != nil {
			log.Warning(e.Error())
		}
		result = append(result, t)
	}

	return result
}
