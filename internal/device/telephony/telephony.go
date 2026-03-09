package telephony

import (
	"adb-backup/internal/shell"
	"strings"
)

type TelephonyManager interface {
	GetSimState() ([]string, error)
	GetSimOperatorAlpha() ([]string, error)
	GetSimOperatorIso() ([]string, error)
	GetSimOperatorNumeric() ([]string, error)
	GetNetworkType() ([]string, error)
	GetNetworkTypeVisualName() ([]string, error)
	GetOperatorAlpha() ([]string, error)
	GetOperatorNumeric() ([]string, error)
}

func NewTelephonyManager(s shell.Shell) TelephonyManager {
	return &shellTelephony{s: s}
}

type shellTelephony struct {
	s shell.Shell
}

func (t *shellTelephony) GetSimState() ([]string, error) {
	return shell.GetPropGsmSimState(t.s)
}

func (t *shellTelephony) GetSimOperatorAlpha() ([]string, error) {
	return shell.GetPropGsmSimOperatorAlpha(t.s)
}

func (t *shellTelephony) GetSimOperatorIso() ([]string, error) {
	return shell.GetPropGsmSimOperatorIso(t.s)
}

func (t *shellTelephony) GetSimOperatorNumeric() ([]string, error) {
	return shell.GetPropGsmSimOperatorNumeric(t.s)
}

func (t *shellTelephony) GetNetworkType() ([]string, error) {
	return shell.GetPropGsmNetworkType(t.s)
}

func (t *shellTelephony) GetNetworkTypeVisualName() ([]string, error) {
	types, err := shell.GetPropGsmNetworkType(t.s)
	if err != nil {
		return nil, err
	}

	visualNames := make([]string, len(types))
	for i, typeValue := range types {
		visualNames[i] = getGsmNetworkVisualName(typeValue)
	}
	return visualNames, nil
}

func (t *shellTelephony) GetOperatorAlpha() ([]string, error) {
	return shell.GetPropGsmOperatorAlpha(t.s)
}

func (t *shellTelephony) GetOperatorNumeric() ([]string, error) {
	return shell.GetPropGsmOperatorNumeric(t.s)
}

// 将 getprop gsm.network.type 的值转换为手机信号栏显示的文本
func getGsmNetworkVisualName(propValue string) string {
	// 统一转为大写处理，兼容不同厂商的 getprop 输出格式
	networkType := strings.ToUpper(strings.TrimSpace(propValue))

	switch networkType {
	// --- 5G 系列 ---
	case "NR", "NR_SA", "NR_NSA":
		return "5G"

	// --- 4G 系列 ---
	case "LTE":
		return "4G"
	case "LTE_CA", "LTE-A":
		return "4G+"

	// --- 3G 系列 ---
	case "HSPAP", "HSPA+":
		return "H+"
	case "HSDPA", "HSUPA", "HSPA":
		return "H"
	case "UMTS", "TD-SCDMA", "WCDMA", "EVDO_0", "EVDO_A", "EVDO_B":
		return "3G"

	// --- 2G 系列 ---
	case "EDGE":
		return "E"
	case "GPRS":
		return "G"
	case "GSM", "CDMA", "1XRTT", "IDEN":
		return "2G"

	// --- 特殊情况 ---
	case "UNKNOWN", "":
		return "No Service"
	default:
		return networkType // 如果是无法识别的新类型，直接返回原值
	}
}
