package shell

import (
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
	return strings.TrimSpace(command), nil
}

func SettingGetBool(d *adb.Device, sys string) (bool, error) {
	command, err := SettingsGet(d, sys)
	if err != nil {
		return false, err
	}
	return command == "1" || command == "true", nil
}

func SettingGetInt(d *adb.Device, sys string) (int, error) {
	command, err := SettingsGet(d, sys)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(command))
}

func SettingGetWifiOn(d *adb.Device) (bool, error) {
	return SettingGetBool(d, "global wifi_on")
}

func SettingsGetAndroidId(d *adb.Device) (string, error) {
	return SettingsGet(d, "secure android_id")
}
