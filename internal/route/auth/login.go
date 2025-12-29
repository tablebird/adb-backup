package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "login", gin.H{})
	}
}
