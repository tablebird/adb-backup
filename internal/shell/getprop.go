package shell

import (
	"strconv"
	"strings"
)

const (
	_GET_PROP = "getprop"
)

func GetProp(s Shell, prop string) (string, error) {
	command, err := s.RunCommand(_GET_PROP + " " + prop)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(command), nil
}

func GetPropInt(s Shell, prop string) (int, error) {
	command, err := GetProp(s, prop)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(command))
}

func GetPropCommaArray(s Shell, prop string) ([]string, error) {
	command, err := GetProp(s, prop)
	if err != nil {
		return nil, err
	}
	split := strings.Split(command, ",")
	for i := range split {
		split[i] = strings.TrimSpace(split[i])
	}
	return split, nil
}

func GetPropGsmOperatorAlpha(s Shell) ([]string, error) {
	return GetPropCommaArray(s, "gsm.operator.alpha")
}

func GetPropGsmOperatorNumeric(s Shell) ([]string, error) {
	return GetPropCommaArray(s, "gsm.operator.numeric")
}

func GetPropGsmNetworkType(s Shell) ([]string, error) {
	return GetPropCommaArray(s, "gsm.network.type")
}

// gsm.sim.state=READY,NOT_READY
func GetPropGsmSimState(s Shell) ([]string, error) {
	return GetPropCommaArray(s, "gsm.sim.state")
}

func GetPropGsmSimOperatorAlpha(s Shell) ([]string, error) {
	return GetPropCommaArray(s, "gsm.sim.operator.alpha")
}

// gsm.sim.operator.numeric=112233
func GetPropGsmSimOperatorNumeric(s Shell) ([]string, error) {
	return GetPropCommaArray(s, "gsm.sim.operator.numeric")
}

func GetPropGsmSimOperatorIso(s Shell) ([]string, error) {
	return GetPropCommaArray(s, "gsm.sim.operator.iso-country")
}

// xiaomi persist.radio.active.slots=2
func GetPropPersistRadioActiveSlots(s Shell) (int, error) {
	return GetPropInt(s, "persist.radio.active.slots")
}

// xiaomi ro.telephony.sim_slots.count=2
func GetPropTelephonySimSlotsCount(s Shell) (int, error) {
	return GetPropInt(s, "ro.telephony.sim_slots.count")
}

// xiaomi ro.telephony.default_cdma_sub=0
func GetPropTelephonyDefaultCdmaSub(s Shell) (int, error) {
	return GetPropInt(s, "ro.telephony.default_cdma_sub")
}

func GetPropBuildVersionRelease(s Shell) (int, error) {
	return GetPropInt(s, "ro.build.version.release")
}

func GetPropProductManufacturer(s Shell) (string, error) {
	return GetProp(s, "ro.product.manufacturer")
}

func GetPropProductModel(s Shell) (string, error) {
	return GetProp(s, "ro.product.model")
}

func GetPropProductBrand(s Shell) (string, error) {
	return GetProp(s, "ro.product.brand")
}

func GetPropConfigMarketingName(s Shell) (string, error) {
	return GetProp(s, "ro.config.marketing_name")
}

func GetPropKernelQemu(s Shell) (string, error) {
	return GetProp(s, "ro.kernel.qemu")
}

func GetPropBootQemuAvdName(s Shell) (string, error) {
	return GetProp(s, "ro.boot.qemu.avd_name")
}
