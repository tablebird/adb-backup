package device

import (
	"adb-backup/internal/database"
	dev "adb-backup/internal/device"
	"adb-backup/internal/log"
	"adb-backup/internal/shell"
	"adb-backup/internal/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	adb "github.com/zach-klippenstein/goadb"
)

type DeviceInfo struct {
	Id             string
	Name           string
	Index          int
	Status         string
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
		connectDevices := dev.GetConnectDevices()
		syncing := dev.GetSyncing()

		h := gin.H{
			"OnlineCount": len(connectDevices),
		}
		devices, err := database.FindAllDevices()
		if err != nil {

			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		allDevices := make(map[string]bool)
		log.DebugF("获取设备列表成功，设备数量： %d", len(connectDevices))
		var OfflineCount = 0
		var devicesInfo []DeviceInfo
		for _, device := range devices {
			var deviceInfo = DeviceInfo{
				Id:    device.Id,
				Name:  device.BuildName(),
				Index: 0,
			}

			serial := device.Serial
			allDevices[serial] = true
			wiFiConnected := false

			adbDevice := connectDevices[serial]
			if adbDevice != nil {
				var androidId string
				state, stateErr := adbDevice.State()

				if state == adb.StateOnline {
					androidId, _ = shell.SettingsGetAndroidId(adbDevice)
					if androidId != deviceInfo.Id {
						OfflineCount++
						devicesInfo = append(devicesInfo, deviceInfo)
						continue
					}
					if _, ok := syncing[serial]; ok {
						deviceInfo.Sync = true
					}
					simStates, _ := shell.GetPropGsmSimState(adbDevice)
					simCount := len(simStates)
					if simCount > 0 {
						var sims []SimInfo
						simAlphas, _ := shell.GetPropGsmSimOperatorAlpha(adbDevice)
						isos, _ := shell.GetPropGsmSimOperatorIso(adbDevice)
						networkTypes, _ := shell.GetPropGsmNetworkType(adbDevice)
						alphas, _ := shell.GetPropGsmOperatorAlpha(adbDevice)
						numerics, _ := shell.GetPropGsmOperatorNumeric(adbDevice)
						SimNumerics, _ := shell.GetPropGsmSimOperatorNumeric(adbDevice)
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
								NetworkType: shell.GetGsmNetworkVisualName(networkType),
								Iso:         isos[i],
							})
						}
						deviceInfo.Sims = sims
					}

					wifi, _ := shell.SettingsGetWifiOn(adbDevice)
					if wifi {
						wifiSSID, err := shell.DumpWifiInfoSsid(adbDevice)
						deviceInfo.WiFiSSID = wifiSSID
						wiFiConnected = err == nil
						deviceInfo.WiFiConnected = wiFiConnected
					}
					batteryLevel, _ := shell.DumpBatteryLevel(adbDevice)
					deviceInfo.Battery = batteryLevel
					poweredTypes, _ := shell.DumpBatteryPoweredType(adbDevice)
					deviceInfo.Charging = len(poweredTypes) > 0
					androidVersion, _ := shell.GetPropBuildVersionRelease(adbDevice)
					deviceInfo.AndroidVersion = androidVersion
				}
				if stateErr == nil {
					deviceInfo.Status = state.String()
				}
			} else {
				OfflineCount++
			}
			devicesInfo = append(devicesInfo, deviceInfo)
		}
		//没用授权通过得设备不会入库
		for _, device := range connectDevices {
			info, err := device.DeviceInfo()
			if err != nil {
				break
			}
			if _, ok := allDevices[info.Serial]; !ok {
				OfflineCount++
			}
			allDevices[info.Serial] = true
		}
		h["Devices"] = devicesInfo
		h["OfflineCount"] = OfflineCount
		h["TotalCount"] = len(allDevices)
		h["LastRefreshTime"] = time.Now().Format("2006-01-02 15:04:05")

		c.HTML(http.StatusOK, "device", h)
	}
}
