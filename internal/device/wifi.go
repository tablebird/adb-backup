package device

import (
	"adb-backup/internal/shell"
)

type WifiManager interface {
	IsEnabled() bool

	WifiSSid() (ssid string, connect bool)
}

type shellWifi struct {
	s shell.Shell
}

func (w *shellWifi) IsEnabled() bool {
	wifi, err := shell.SettingsGetWifiOn(w.s)
	if err != nil {
		return false
	}
	return wifi
}

func (w *shellWifi) WifiSSid() (ssid string, connect bool) {
	wifiSSID, err := shell.DumpWifiInfoSsid(w.s)
	ssid = wifiSSID
	connect = err == nil
	return
}
