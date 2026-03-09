package validator

import (
	"adb-backup/internal/database"
	"adb-backup/internal/device"
	"adb-backup/internal/web/base"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func deviceId(fl validator.FieldLevel) bool {
	dev := getDevice(fl)
	return dev != nil
}

func getDevice(fl validator.FieldLevel) *database.Device {
	deviceId, ok := fl.Field().Interface().(string)
	if ok && deviceId != "" {
		dbDevice, err := database.FindDeviceById(deviceId)
		if err != nil {
			return nil
		}
		ctx := getContext(fl)
		if ctx != nil {
			ctx.Set(base.ContextDeviceIdKey, dbDevice)
			ctx.SetTypeKey(dbDevice)
		}
		return &dbDevice
	} else {
		return nil
	}
}

func deviceIdConnect(fl validator.FieldLevel) bool {
	dev := getDevice(fl)
	if dev != nil {
		adbDevice := device.GetDevice(dev.Serial)
		if adbDevice != nil {
			ctx := getContext(fl)
			if ctx != nil {
				ctx.SetTypeKey(adbDevice)
			}
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
		dbDevice, err := database.FindDeviceById(deviceId)
		if err != nil {
			base.RespJsonBadRequest(c, "设备不存在")
			c.Abort()
			return
		}
		base.SetContextTypeKey(c, dbDevice)
	}
}

func DeviceIdConnectMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		DeviceIdMiddleware()(c)
		if c.IsAborted() {
			return
		}
		dbDevice := c.MustGet(base.TypeKey[database.Device]()).(database.Device)

		adbDevice := device.GetDevice(dbDevice.Serial)

		if adbDevice == nil {
			base.RespJsonBadRequest(c, "设备未连接")
			c.Abort()
			return
		}
		base.SetContextTypeKey(c, adbDevice)
	}
}
