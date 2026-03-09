package shell

import (
	"strings"
)

func GetMarketingName(s Shell) (string, error) {
	qemu, _ := GetPropKernelQemu(s)
	if qemu == "1" {
		avdName, _ := GetPropBootQemuAvdName(s)
		if avdName != "" {
			return strings.ReplaceAll(avdName, "_", " "), nil
		}
	}
	manufacturer, err := GetPropProductManufacturer(s)
	if err != nil {
		return "", err
	}
	var marketingName = ""
	switch manufacturer {
	case "HUAWEI":
		marketingName, _ = GetPropConfigMarketingName(s)
	}
	if marketingName != "" {
		return marketingName, nil
	}
	return GetPropProductModel(s)
}
