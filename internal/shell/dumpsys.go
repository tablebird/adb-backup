package shell

import (
	"errors"
	"strconv"
	"strings"

	adb "github.com/zach-klippenstein/goadb"
)

const (
	_DUMP_SYS = "dumpsys"
)

func DumpSys(d *adb.Device, sys string) (string, error) {
	command, err := d.RunCommand(_DUMP_SYS + " " + sys)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(command), nil
}

func DumpWifiInfoSsid(d *adb.Device) (string, error) {
	res, err := DumpSys(d, "wifi  | grep \"mWifiInfo SSID:\"")
	if err != nil {
		return "", err
	}

	return _parseWifiInfoSsid(res)
}

func _parseWifiInfoSsid(res string) (string, error) {
	lines := strings.Split(res, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "mWifiInfo SSID:") {
			split := strings.Split(line, ",")
			ssid := strings.TrimSpace(split[0][len("mWifiInfo SSID:"):])
			if ssid == "<unknown ssid>" {
				return "", errors.New("not connect")
			}
			return ssid[1 : len(ssid)-1], nil
		}
	}
	return "", errors.New("not connect")
}

func DumpBatteryLevel(d *adb.Device) (int, error) {
	res, err := DumpSys(d, "battery | grep level:")
	if err != nil {
		return int(0), err
	}
	return _parseBatteryLevel(res)
}

func _parseBatteryLevel(res string) (int, error) {
	lines := strings.Split(res, "\n")
	var level int
	var err error
	for _, line := range lines {
		// 解析 "level: 80" 格式的输出
		parts := strings.Split(strings.TrimSpace(line), ":")
		if len(parts) == 2 {
			level, err = strconv.Atoi(strings.TrimSpace(parts[1]))
			if err == nil {
				return level, nil
			}
		}
	}
	if err == nil {
		err = errors.New("battery level not found")
	}
	return int(0), err
}

func DumpBatteryPoweredType(d *adb.Device) ([]BatteryPoweredType, error) {
	res, err := DumpSys(d, "battery | grep powered:")
	if err != nil {
		return nil, err
	}
	return _parseBatteryPoweredType(res), nil
}

func _parseBatteryPoweredType(res string) []BatteryPoweredType {
	var powereds []BatteryPoweredType
	lines := strings.Split(res, "\n")
	for _, line := range lines {
		// 解析 "powered: AC" 格式的输出
		parts := strings.Split(strings.TrimSpace(line), ":")
		if len(parts) == 2 {
			value := strings.TrimSpace(parts[1])
			if value == "true" {
				key := strings.TrimSpace(parts[0])
				powerName := strings.ReplaceAll(key, "powered", "")
				poweredType, err := parseBatteryPoweredType(strings.TrimSpace(powerName))
				if err == nil {
					powereds = append(powereds, poweredType)
				}
			}
		}
	}
	return powereds
}
