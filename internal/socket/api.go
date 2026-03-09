package socket

import (
	"adb-backup/internal/web/base"

	"github.com/gin-gonic/gin"
)

func Socket(c *gin.Context) {
	client, err := NewClient(c)
	if err != nil {
		base.RespJsonInternalServerError(c, "websocket error")
		return
	}
	defer client.Close()

	addClient(client)

	client.doWork()
}
