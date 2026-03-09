package sms

import (
	"adb-backup/internal/config"
	"adb-backup/internal/database"
	"adb-backup/internal/device"
	"adb-backup/internal/shell"
	"adb-backup/internal/web/base"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SmsPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceId := c.GetString(base.ContextDeviceIdKey)
		h := gin.H{
			"DeviceID": deviceId,
		}
		dbDevice := c.MustGet(base.TypeKey[database.Device]()).(database.Device)

		h["DeviceName"] = dbDevice.BuildName()
		if config.Feature.EnableSendSms {
			adbDevice := device.GetDevice(dbDevice.Serial)
			if adbDevice != nil {
				networkTypes, err := shell.GetPropGsmNetworkType(adbDevice)
				if err == nil {
					h["NetworkTypes"] = networkTypes
					var enableSendSms bool
					for _, networkType := range networkTypes {
						if networkType != "" && strings.ToUpper(networkType) != "UNKNOWN" {
							enableSendSms = true
							break
						}
					}
					h["EnableSendSms"] = enableSendSms
				}
			}
		}
		c.HTML(http.StatusOK, "sms_detail", h)
	}
}
