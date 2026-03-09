package login

import (
	"adb-backup/internal/config"
	"adb-backup/internal/database"
	"adb-backup/internal/web/base"

	"github.com/gin-gonic/gin"
)

func Login() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		var req struct {
			Username   string `json:"username"`
			Password   string `json:"password"`
			AuthSource int    `json:"auth_source"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			base.RespJsonBadRequest(c, "请求参数错误")
			return
		}

		user, err := database.Authenticate(req.AuthSource, req.Username, req.Password)
		if err != nil {
			base.RespJsonBadRequest(c, err.Error())
			return
		}

		token := generateToken()
		tokenMap[token] = user.Id
		c.SetCookie("login_token",
			token,                   // 令牌值
			config.Web.CookieMaxAge, // 有效期（秒）
			"/",                     // 作用路径
			"",                      // 作用域名
			false,                   // Secure（HTTPS下启用）
			true,                    // HttpOnly
		)
		base.RespJsonSuccess(c, "登录成功", nil)
	})
}

func LoginOut() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.SetCookie("login_token", "", -1, "/", "", false, true)
		base.RespJsonSuccess(c, "登出成功", nil)
	})
}

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		base.RespJsonSuccess(c, "", nil)
	}
}
