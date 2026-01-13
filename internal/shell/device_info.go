package shell

import (
	"strings"

	adb "github.com/zach-klippenstein/goadb"
)

func GetMarketingName(d *adb.Device) (string, error) {
	qemu, _ := GetPropKernelQemu(d)
	if qemu == "1" {
		avdName, _ := GetPropBootQemuAvdName(d)
		if avdName != "" {
			return strings.ReplaceAll(avdName, "_", " "), nil
		}
	}
	manufacturer, err := GetPropProductManufacturer(d)
	if err != nil {
		return "", err
	}
	var marketingName = ""
	switch manufacturer {
	case "HUAWEI":
		marketingName, _ = GetPropConfigMarketingName(d)
	}
	if marketingName != "" {
		return marketingName, nil
	}
	return GetPropProductModel(d)
}
