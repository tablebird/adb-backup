package login

import (
	"adb-backup/internal/config"
	"adb-backup/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		var req struct {
			Username   string `json:"username"`
			Password   string `json:"password"`
			AuthSource string `json:"auth_source"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		user, err := database.FindUserByName(req.Username)
		if err != nil || !user.CheckPassword(req.Password, user.Password) {
			c.JSON(401, gin.H{"error": "认证失败"})
			return
		}
		token := generateToken()
		tokenMap[token] = user.Id
		c.SetCookie("login_token",
			token,                   // 令牌值
			config.Web.CookieMaxAge, // 有效期（秒）
			"/",                     // 作用路径
			"",                      // 作用域名
			true,                    // Secure（HTTPS下启用）
			true,                    // HttpOnly
		)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "登录成功",
		})
	})
}

func LoginOut() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.SetCookie("login_token", "", -1, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "登出成功",
		})
	})
}
