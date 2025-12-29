package main

import (
	"adb-backup/internal/admin"
	"adb-backup/internal/config"
	"adb-backup/internal/database"
	"adb-backup/internal/device"
	"adb-backup/internal/log"
	"adb-backup/internal/notify"
	"adb-backup/internal/web"
)

func main() {
	log.InfoF("服务启动中....")

	config.Conf.InitConfig()
	url := config.Conf.NotifyWebhookUrl
	if len(url) != 0 {
		notify.Notify = notify.Webhook{
			Url: url,
		}
	}
	database.InitDB()
	admin.InitAdmin()
	go web.InitWeb()
	device.StartWatch()
}
