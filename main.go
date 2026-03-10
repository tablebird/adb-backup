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

	config.InitConfig()
	notify.InitNotify()
	database.InitDB()
	admin.InitAdmin()
	go web.InitWeb()
	device.StartWatch()
}
