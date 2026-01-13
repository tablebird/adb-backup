package sms

import (
	"adb-backup/internal/config"
	"adb-backup/internal/database"
	"adb-backup/internal/device"
	"adb-backup/internal/shell"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SmsPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		deviceId := r.URL.Query().Get("device_id")
		if deviceId == "" {
			c.Redirect(http.StatusFound, "/")
			return
		}
		h := gin.H{
			"DeviceID": deviceId,
		}
		dbDevice, err := database.FindDeviceById(deviceId)
		if err == nil && dbDevice.Id != "" {
			h["DeviceName"] = dbDevice.BuildName()
		}
		if config.Feature.EnableSendSms {
			adbDevice := device.GetDevice(deviceId)
			if adbDevice != nil {
				networkTypes, err := shell.GetPropGsmNetworkType(adbDevice)
				if err == nil {
					h["NetworkTypes"] = networkTypes
				}
			}
		}
		c.HTML(http.StatusOK, "sms_detail", h)
	}
}
