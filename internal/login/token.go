package login

import (
	"adb-backup/internal/web/base"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	tokenMap = make(map[string]int) // token到用户ID的映射
)

func NotAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查Cookie中的登录态
		token, err := c.Cookie("login_token")
		if err != nil || token == "" {
			// 未登录，继续执行
			c.Next()
			return
		}

		_, validToken := tokenMap[token]
		if !validToken {
			// token无效，清除Cookie并继续执行
			c.SetCookie("login_token", "", -1, "/", "", false, true)
			c.Next()
			return
		}

		// 登录态有效，跳过后续中间件
		c.Redirect(http.StatusFound, "/")
		c.Abort()
	}
}

func getLoginUrl(c *gin.Context) string {
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	if query != "" {
		path += "?" + query
	}
	redirect := url.QueryEscape(path)
	if (redirect == "/") || (redirect == "") {
		return "/login"
	}
	return "/login?redirect=" + redirect
}

func AuthMiddleware() gin.HandlerFunc {
	return authMiddleware(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api") {
			base.RespJson(c, http.StatusUnauthorized, "请先登录", nil)
		} else if "/ws" == path {
			c.String(http.StatusUnauthorized, "Unauthorized: Please login first")
		} else {
			c.Redirect(http.StatusFound, getLoginUrl(c))
		}
	})
}

func authMiddleware(failFunc func(*gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查Cookie中的登录态
		token, err := c.Cookie("login_token")
		if err != nil || token == "" {

			// 未登录，跳转到登录页
			failFunc(c)
			c.Abort()
			return
		}

		userId, validToken := tokenMap[token]
		if !validToken {

			// token无效，清除Cookie并跳转登录
			c.SetCookie("login_token", "", -1, "/", "", false, true)
			failFunc(c)
			c.Abort()
			return
		}
		c.Set("userId", userId)

		// 登录态有效，继续执行
		c.Next()
	}
}

func generateToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
