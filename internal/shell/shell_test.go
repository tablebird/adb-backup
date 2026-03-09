package shell

import adb "github.com/zach-klippenstein/goadb"

const (
	DEVICE_TEST_SERIAL = "emulator-5554"
)

func BuildShell() Shell {
	client, _ := adb.NewWithConfig(adb.ServerConfig{})
	device := client.Device(adb.DeviceWithSerial(DEVICE_TEST_SERIAL))
	return device
}
