package device

import (
	"adb-backup/internal/shell"

	adb "github.com/zach-klippenstein/goadb"
)

type WifiManager interface {
	IsEnabled() bool

	WifiSSid() (ssid string, connect bool)
}

type shellWifi struct {
	adbDevice *adb.Device
}

func (w *shellWifi) IsEnabled() bool {
	wifi, err := shell.SettingsGetWifiOn(w.adbDevice)
	if err != nil {
		return false
	}
	return wifi
}

func (w *shellWifi) WifiSSid() (ssid string, connect bool) {
	wifiSSID, err := shell.DumpWifiInfoSsid(w.adbDevice)
	ssid = wifiSSID
	connect = err == nil
	return
}
