package device

import adb "github.com/tablebird/goadb"

type DeviceState int8

const (
	StateInvalid DeviceState = iota
	StateUnauthorized
	StateDisconnected
	StateOffline
	StateOnline
	StateError
)

var deviceStateStrings = map[adb.DeviceState]DeviceState{
	adb.StateInvalid:      StateInvalid,
	adb.StateDisconnected: StateDisconnected,
	adb.StateOffline:      StateOffline,
	adb.StateOnline:       StateOnline,
	adb.StateUnauthorized: StateUnauthorized,
}

func (s DeviceState) String() string {
	switch s {
	case StateInvalid:
		return "Invalid"
	case StateDisconnected:
		return "Disconnected"
	case StateOffline:
		return "Offline"
	case StateOnline:
		return "Online"
	case StateUnauthorized:
		return "Unauthorized"
	default:
		return "Error"
	}
}

func deviceStateToStr(state adb.DeviceState) DeviceState {
	return deviceStateStrings[state]
}
