package power

import (
	"fmt"
)

type BatteryPoweredType int8

const (
	BatteryPoweredTypeUnknown BatteryPoweredType = iota
	BatteryPoweredTypeAC
	BatteryPoweredTypeUSB
	BatteryPoweredTypeWireless
)

var (
	batteryPoweredTypeNames = map[string]BatteryPoweredType{
		"Unknown":  BatteryPoweredTypeUnknown,
		"AC":       BatteryPoweredTypeAC,
		"USB":      BatteryPoweredTypeUSB,
		"Wireless": BatteryPoweredTypeWireless,
	}
)

func parseBatteryPoweredType(str string) (BatteryPoweredType, error) {
	poweredType, ok := batteryPoweredTypeNames[str]
	if ok {
		return poweredType, nil
	}
	return BatteryPoweredTypeUnknown, fmt.Errorf("invalid battery powered type: %q", str)
}

func (b BatteryPoweredType) String() string {
	for name, poweredType := range batteryPoweredTypeNames {
		if poweredType == b {
			return name
		}
	}
	return fmt.Sprintf("BatteryPoweredType(%d)", b)
}
