package shell

import (
	"errors"
	"strconv"
	"strings"
)

const (
	_GETEVENT = "getevent"
)

func GetEvent(s Shell, args ...string) (string, error) {
	return s.RunCommand(_GETEVENT, args...)
}

func GetEventTouch(s Shell) (string, int, int, error) {
	res, err := GetEvent(s, "-il")
	if err != nil {
		return "", -1, -1, err
	}
	devices := strings.Split(res, "add device ")
	for _, device := range devices {
		lines := strings.Split(device, "\n")
		var x, y int
		var device string
		for _, line := range lines {
			if strings.Contains(line, "/dev/input/event") {
				device = strings.SplitN(line, ":", 2)[1]
			} else if strings.HasPrefix(line, "name") {
				if !strings.Contains(strings.ToUpper(line), "TOUCH") {
					break
				}
			} else {
				if strings.Contains(line, "ABS_MT_POSITION_X") {
					x = _parseEventMax(line)
				} else if strings.Contains(line, "ABS_MT_POSITION_Y") {
					y = _parseEventMax(line)
				}
			}
		}
		if device != "" && x > 0 && y > 0 {
			return strings.TrimSpace(device), x, y, nil
		}
	}
	return res, -1, -1, errors.New("not found touch device")
}

func _parseEventMax(line string) int {
	values := strings.Split(strings.SplitN(line, ":", 2)[1], ", ")
	for _, item := range values {
		if strings.HasPrefix(item, "max ") {
			max, err := strconv.Atoi(strings.ReplaceAll(item, "max ", ""))
			if err != nil {
				return -1
			}
			return max
		}
	}
	return -1
}
