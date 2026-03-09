package device

import (
	"adb-backup/internal/device"
	"adb-backup/internal/web/base"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeviceDetail(c *gin.Context) {
	h := gin.H{
		"DeviceId": c.MustGet(base.ContextDeviceIdKey),
	}
	dev := c.MustGet(base.ContextDeviceKey).(device.ConnectDevice)
	display := dev.GetDisplay()
	if display != nil {
		width, height := display.Size()
		h["Width"] = width
		h["Height"] = height
		h["DisplayState"] = display.IsOn()
	}
	touch := dev.GetTouch()
	width, height := touch.Size()
	if touch != nil {
		h["TouchWidth"] = width
		h["TouchHeight"] = height
	}

	c.HTML(http.StatusOK, "deviceDetail", h)
}
