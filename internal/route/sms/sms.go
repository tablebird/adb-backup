package sms

import (
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
		c.HTML(http.StatusOK, "sms_detail", h)
	}
}
