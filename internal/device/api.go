package device

import (
	"adb-backup/internal/web/base"

	"github.com/gin-gonic/gin"
)

func RefreshScanDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		checkAndSyncDevices()
		base.RespJsonSuccess(c, "刷新成功", nil)
	}
}
