package auth

import (
	"adb-backup/internal/auth"
	"adb-backup/internal/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := gin.H{}

		keys := make([]auth.Type, 0, len(auth.Providers))
		for key, value := range auth.Providers {
			if !value.IsReady() {
				continue
			}
			keys = append(keys, key)
		}
		log.DebugF("AuthSource %v", keys)
		h["AuthSources"] = keys
		appendUser(h)

		c.HTML(http.StatusOK, "login", h)
	}
}
