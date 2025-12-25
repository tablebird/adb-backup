package web

import (
	"adb-backup/internal/config"
	deviceApi "adb-backup/internal/device"
	"adb-backup/internal/log"
	"adb-backup/internal/route/device"
	"adb-backup/internal/route/sms"
	smsApi "adb-backup/internal/sms"
	"adb-backup/internal/utils"
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

	r.GET("/", device.DevicesInfo())
	r.GET("/api/device/refreshScan", deviceApi.RefreshScanDevice())
	r.GET("/sms", sms.SmsPage())
	r.GET("/api/sms/conversations", smsApi.GetConversationsApiHandler())
	r.GET("/api/sms/messages/latest", smsApi.GetLatestMessagesApiHandler())
	r.GET("/api/sms/messages/old", smsApi.GetOldMessagesApiHandler())
	port := config.Conf.WebPort
	logWebUrl(port)
	r.Run(fmt.Sprintf(":%d", port))
}

func logWebUrl(port int) {
	ip, err := utils.GetLocalHostIP()
	if err == nil {
		log.InfoF("web url: http://%s:%d", ip, port)
	} else {
		log.ErrorF("get localhost ip error: %v", err)
	}

}
