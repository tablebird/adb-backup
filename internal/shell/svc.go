package shell

import (
	"fmt"
)

const (
	_SVC = "svc"
)

func Svc(s Shell, method string, args ...string) (string, error) {
	res, err := s.RunCommand(_SVC, append([]string{method}, args...)...)
	if err != nil {
		return "", err
	}
	return res, nil
}

func SvcEnable(s Shell, method string, enable bool) (bool, error) {
	var arg string
	if enable {
		arg = "enable"
	} else {
		arg = "disable"
	}
	res, err := Svc(s, method, arg)
	if err != nil {
		return false, err
	}
	if res == "" {
		return true, nil
	}
	return false, fmt.Errorf("%s", res)
}

func SvcWifi(s Shell, enable bool) (bool, error) {
	return SvcEnable(s, "wifi", enable)
}

func SvcData(s Shell, enable bool) (bool, error) {
	return SvcEnable(s, "data", enable)
}
