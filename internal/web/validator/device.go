package validator

import (
	"adb-backup/internal/device"
	"adb-backup/internal/web/base"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func deviceId(fl validator.FieldLevel) bool {
	dev := getDevice(fl)
	return dev != nil
}

func getDevice(fl validator.FieldLevel) device.Device {
	deviceId, ok := fl.Field().Interface().(string)
	if ok && deviceId != "" {
		dev, err := device.FindDeviceById(deviceId)
		if err != nil {
			return nil
		}
		ctx := getContext(fl)
		if ctx != nil {
			ctx.Set(base.ContextDeviceIdKey, deviceId)
			ctx.Set(base.ContextDeviceKey, dev)
		}
		return dev
	} else {
		return nil
	}
}

func deviceIdConnect(fl validator.FieldLevel) bool {
	dev := getDevice(fl)
	if dev != nil {
		if _, ok := dev.(device.ConnectDevice); ok {
			return true
		}
	}
	return false
}

func DeviceIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceId := c.Request.URL.Query().Get("device_id")
		if deviceId == "" {
			base.RespJsonBadRequest(c, "设备ID不能为空")
			c.Abort()
			return
		}
		c.Set(base.ContextDeviceIdKey, deviceId)
		dev, err := device.FindDeviceById(deviceId)
		if err != nil {
			base.RespJsonBadRequest(c, "设备不存在")
			c.Abort()
			return
		}
		c.Set(base.ContextDeviceKey, dev)
	}
}

func DeviceIdConnectMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		DeviceIdMiddleware()(c)
		if c.IsAborted() {
			return
		}
		dev := c.MustGet(base.ContextDeviceKey).(device.Device)
		if _, ok := dev.(device.ConnectDevice); !ok {
			base.RespJsonBadRequest(c, "设备未连接")
			c.Abort()
		}
	}
}
