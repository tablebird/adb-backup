package device

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RefreshScanDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		checkAndSyncDevices()
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "刷新成功",
			"data": nil,
		})
	}
}
