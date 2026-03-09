package base

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BuildResp(code int, msg string, data any) any {
	return gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	}
}

func RespJson(c *gin.Context, code int, msg string, data any) {
	c.JSON(code, BuildResp(code, msg, data))
}

func RespJsonSuccess(c *gin.Context, msg string, data any) {
	RespJson(c, http.StatusOK, msg, data)
}

func RespJsonBadRequest(c *gin.Context, msg string) {
	RespJson(c, http.StatusBadRequest, msg, nil)
}

func RespJsonInternalServerError(c *gin.Context, msg string) {
	RespJson(c, http.StatusInternalServerError, msg, nil)
}
