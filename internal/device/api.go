package device

import (
	"adb-backup/internal/database"
	"adb-backup/internal/log"
	"adb-backup/internal/web/base"

	"github.com/gin-gonic/gin"
)

type EnableStatusNotifyReq struct {
	base.ContextReq
	DeviceId     string `json:"device_id" binding:"required,deviceId"`
	StatusNotify *bool  `json:"status_notify" binding:"required"`
}

func RefreshScanDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		scanAllDevices()
		base.RespJsonSuccess(c, "刷新成功", nil)
	}
}

func EnableStatusNotify() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req EnableStatusNotifyReq
		if err := c.ShouldBindJSON(&req); err != nil {
			base.RespJsonBadRequest(c, "参数错误")
			return
		}
		log.DebugF("EnableStatusNotify deviceId %s notify %t", req.DeviceId, *req.StatusNotify)
		err := database.UpdateDeviceStatusNotify(req.DeviceId, *req.StatusNotify)
		if err != nil {
			base.RespJsonInternalServerError(c, "更新中状态通知失败")
			return
		}
		base.RespJsonSuccess(c, "更新状态通知成功", nil)
	}
}
