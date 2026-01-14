package shell

import (
	"fmt"
	"strconv"
	"strings"

	adb "github.com/zach-klippenstein/goadb"
)

const (
	_SETTINGS_GET = "settings get"
)

func SettingsGet(d *adb.Device, sys string) (string, error) {
	command, err := d.RunCommand(_SETTINGS_GET + " " + sys)
	if err != nil {
		return command, err
	}
	if strings.Contains(command, "Can't find service: settings") {
		return "", fmt.Errorf("%s", command)
	}
	return strings.TrimSpace(command), nil
}

func SettingsGetBool(d *adb.Device, sys string) (bool, error) {
	command, err := SettingsGet(d, sys)
	if err != nil {
		return false, err
	}
	return command == "1" || command == "true", nil
}

func SettingsGetInt(d *adb.Device, sys string) (int, error) {
	command, err := SettingsGet(d, sys)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(command))
}

func SettingsGetWifiOn(d *adb.Device) (bool, error) {
	return SettingsGetBool(d, "global wifi_on")
}

func SettingsGetAndroidId(d *adb.Device) (string, error) {
	return SettingsGet(d, "secure android_id")
}
