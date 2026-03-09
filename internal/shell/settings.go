package shell

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	_SETTINGS_GET = "settings get"
)

func SettingsGet(s Shell, sys string) (string, error) {
	command, err := s.RunCommand(_SETTINGS_GET + " " + sys)
	if err != nil {
		return command, err
	}
	if strings.Contains(command, "Can't find service: settings") {
		return "", fmt.Errorf("%s", command)
	}
	return strings.TrimSpace(command), nil
}

func SettingsGetBool(s Shell, sys string) (bool, error) {
	command, err := SettingsGet(s, sys)
	if err != nil {
		return false, err
	}
	return command == "1" || command == "true", nil
}

func SettingsGetInt(s Shell, sys string) (int, error) {
	command, err := SettingsGet(s, sys)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(command))
}

func SettingsGetWifiOn(s Shell) (bool, error) {
	return SettingsGetBool(s, "global wifi_on")
}

func SettingsGetAndroidId(s Shell) (string, error) {
	return SettingsGet(s, "secure android_id")
}
