package sms

import (
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
		device := device.GetDevice(deviceId)
		if device != nil {
			networkTypes, err := shell.GetPropGsmNetworkType(device)
			if err == nil {
				h["NetworkTypes"] = networkTypes
			}
		}
		c.HTML(http.StatusOK, "sms_detail", h)
	}
}
