package device

import (
	"adb-backup/internal/device"
	"adb-backup/internal/log"
	"adb-backup/internal/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DeviceInfo struct {
	Id             string
	Name           string
	Index          int
	Status         string
	StatusNotify   bool
	Sync           bool
	Sims           []SimInfo
	WiFiConnected  bool
	WiFiSSID       string
	Battery        int
	Charging       bool
	AndroidVersion int
}

type SimInfo struct {
	Operator    string
	NetworkType string
	Iso         string
}

func DevicesInfo() gin.HandlerFunc {
	return func(c *gin.Context) {

		devices := device.FindAllDevices()
		total := len(devices)
		log.DebugF("获取设备列表成功，设备数量： %d", total)

		var onlineCount = 0
		var deviceInfos []DeviceInfo
		for i, item := range devices {
			var deviceInfo = DeviceInfo{
				Id:     item.Id(),
				Name:   item.Name(),
				Index:  i,
				Status: item.State().String(),
			}
			deviceDB := item.GetDeviceDB()
			if deviceDB != nil {
				deviceInfo.StatusNotify = deviceDB.StatusNotify
			}
			if item.State() == device.StateOnline {
				onlineCount++
				if item, ok := item.(device.ConnectDevice); ok {
					sync := item.GetSync()
					deviceInfo.Sync = sync != nil && sync.IsSyncing()
					telephone := item.GetTelephony()
					if telephone != nil {
						simStates, _ := telephone.GetSimState()
						simCount := len(simStates)
						if simCount > 0 {
							var sims []SimInfo
							simAlphas, _ := telephone.GetSimOperatorAlpha()
							isos, _ := telephone.GetSimOperatorIso()
							networkTypes, _ := telephone.GetNetworkTypeVisualName()
							alphas, _ := telephone.GetOperatorAlpha()
							numerics, _ := telephone.GetOperatorNumeric()
							SimNumerics, _ := telephone.GetSimOperatorNumeric()
							for i, simState := range simStates {
								if simState == "ABSENT" {
									continue
								}
								var operator string
								networkType := utils.ArrayGet(networkTypes, i)
								alpha := utils.ArrayGet(alphas, i)
								numeric := utils.ArrayGet(numerics, i)
								SimNumeric := utils.ArrayGet(SimNumerics, i)
								simAlpha := utils.ArrayGet(simAlphas, i)
								if len(alpha) == 0 {
									operator = simAlpha
								} else if numeric == SimNumeric {
									operator = alpha
								} else {
									operator = fmt.Sprintf("%s(%s)", simAlpha, alpha)
								}
								sims = append(sims, SimInfo{
									Operator:    operator,
									NetworkType: networkType,
									Iso:         isos[i],
								})
							}
							deviceInfo.Sims = sims
						}
					}
					wifi := item.GetWifi()
					if wifi != nil && wifi.IsEnabled() {
						ssid, connect := wifi.WifiSSid()
						deviceInfo.WiFiConnected = connect
						deviceInfo.WiFiSSID = ssid
					}
					power := item.GetPower()
					if power != nil {
						deviceInfo.Battery = power.BatteryLevel()
						deviceInfo.Charging = power.CharingType() != nil
					}
					build := item.GetBuild()
					if build != nil {
						deviceInfo.AndroidVersion = build.VersionRelease()
					}
				}
			}
			deviceInfos = append(deviceInfos, deviceInfo)
		}

		h := gin.H{
			"Devices":         deviceInfos,
			"OnlineCount":     onlineCount,
			"TotalCount":      total,
			"OfflineCount":    total - onlineCount,
			"LastRefreshTime": time.Now().Format("2006-01-02 15:04:05"),
		}

		c.HTML(http.StatusOK, "device", h)
	}
}
