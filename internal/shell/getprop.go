package shell

import (
	"strconv"
	"strings"

	adb "github.com/zach-klippenstein/goadb"
)

const (
	_GET_PROP = "getprop"
)

func GetProp(d *adb.Device, prop string) (string, error) {
	command, err := d.RunCommand(_GET_PROP + " " + prop)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(command), nil
}

func GetPropInt(d *adb.Device, prop string) (int, error) {
	command, err := GetProp(d, prop)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(command))
}

func GetPropCommaArray(d *adb.Device, prop string) ([]string, error) {
	command, err := GetProp(d, prop)
	if err != nil {
		return nil, err
	}
	split := strings.Split(command, ",")
	for i := range split {
		split[i] = strings.TrimSpace(split[i])
	}
	return split, nil
}

func GetPropGsmOperatorAlpha(d *adb.Device) ([]string, error) {
	return GetPropCommaArray(d, "gsm.operator.alpha")
}

func GetPropGsmOperatorNumeric(d *adb.Device) ([]string, error) {
	return GetPropCommaArray(d, "gsm.operator.numeric")
}

func GetPropGsmNetworkType(d *adb.Device) ([]string, error) {
	return GetPropCommaArray(d, "gsm.network.type")
}

// gsm.sim.state=READY,NOT_READY
func GetPropGsmSimState(d *adb.Device) ([]string, error) {
	return GetPropCommaArray(d, "gsm.sim.state")
}

func GetPropGsmSimOperatorAlpha(d *adb.Device) ([]string, error) {
	return GetPropCommaArray(d, "gsm.sim.operator.alpha")
}

// gsm.sim.operator.numeric=112233
func GetPropGsmSimOperatorNumeric(d *adb.Device) ([]string, error) {
	return GetPropCommaArray(d, "gsm.sim.operator.numeric")
}

func GetPropGsmSimOperatorIso(d *adb.Device) ([]string, error) {
	return GetPropCommaArray(d, "gsm.sim.operator.iso-country")
}

// xiaomi persist.radio.active.slots=2
func GetPropPersistRadioActiveSlots(d *adb.Device) (int, error) {
	return GetPropInt(d, "persist.radio.active.slots")
}

// xiaomi ro.telephony.sim_slots.count=2
func GetPropTelephonySimSlotsCount(d *adb.Device) (int, error) {
	return GetPropInt(d, "ro.telephony.sim_slots.count")
}

// xiaomi ro.telephony.default_cdma_sub=0
func GetPropTelephonyDefaultCdmaSub(d *adb.Device) (int, error) {
	return GetPropInt(d, "ro.telephony.default_cdma_sub")
}
