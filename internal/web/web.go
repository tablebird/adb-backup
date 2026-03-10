package web

import (
	"adb-backup/internal/config"
	deviceApi "adb-backup/internal/device"
	"adb-backup/internal/log"
	"adb-backup/internal/login"
	"adb-backup/internal/route/auth"
	"adb-backup/internal/route/device"
	"adb-backup/internal/route/sms"
	"adb-backup/internal/screen"
	smsApi "adb-backup/internal/sms"
	"adb-backup/internal/socket"
	"adb-backup/internal/utils"
	"adb-backup/internal/web/validator"
	tmpl "adb-backup/templates"
	"html/template"

	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

type ApiResponse struct {
	Code int         `json:"code"` // 200成功，其他失败
	Msg  string      `json:"msg"`  // 提示信息
	Data interface{} `json:"data"` // 业务数据
}

func InitWeb() {
	log.InfoF("web start")
	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"formatTime": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		// 自定义函数：转大写
		"toUpper": func(s string) string {
			return strings.ToUpper(s)
		},
		"add": func(a, b int) int {
			return a + b
		},
	})
	r.HTMLRender = tmpl.GetHTMLRender(r.FuncMap, render.Delims{Left: "{{", Right: "}}"})

	validator.RegisterValidation()

	r.GET("/login", login.NotAuthMiddleware(), auth.LoginPage())
	r.POST("/api/login", login.Login())
	r.POST("/api/logout", login.LoginOut())
	group := r.Group("/", login.AuthMiddleware())
	group.GET("/api/checkAuth", login.CheckLogin())

	group.GET("/", device.DevicesInfo())
	group.GET("/ws", socket.Socket)
	group.GET("/api/device/refreshScan", deviceApi.RefreshScanDevice())
	deviceIdGroup := group.Group("/", validator.DeviceIdMiddleware())
	deviceIdGroup.GET("/sms", sms.SmsPage())
	deviceIdGroup.GET("/api/sms/conversations", smsApi.GetConversationsApiHandler())
	deviceIdGroup.GET("/api/sms/messages/latest", smsApi.GetLatestMessagesApiHandler())
	deviceIdGroup.GET("/api/sms/messages/old", smsApi.GetOldMessagesApiHandler())
	deviceIdGroup.GET("/api/sms/messages/new", smsApi.GetNewMessageApiHandler())

	deviceConnectGroup := group.Group("/", validator.DeviceIdConnectMiddleware())
	deviceConnectGroup.GET("/device", device.DeviceDetail)
	deviceConnectGroup.GET("/api/screen", screen.ScreenView())
	deviceConnectGroup.GET("/api/screenCap", screen.ScreenCap())

	group.POST("/api/device/statusNotify", deviceApi.EnableStatusNotify())

	group.POST("/api/sms/send", smsApi.SendMessage())

	port := config.Web.WebPort
	address := config.Web.Address
	if address == "" {
		logWebUrl(port)
		r.Run(fmt.Sprintf(":%d", port))
	} else {
		log.InfoF("web url: http://%s:%d", address, port)
		r.Run(fmt.Sprintf("%s:%d", address, port))
	}
}

func logWebUrl(port int) {
	ip, err := utils.GetLocalHostIP()
	if err == nil {
		log.InfoF("web url: http://%s:%d", ip, port)
	} else {
		log.ErrorF("get localhost ip error: %v", err)
	}

}
