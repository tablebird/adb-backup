package database

import (
	"adb-backup/internal/config"
	"fmt"
	"io"
	"log"

	l "adb-backup/internal/log"

	postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func InitDB() {
	var conf = config.DB
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		conf.DbHost, conf.DbPort, conf.DbUser, conf.DbPass, conf.DbName, conf.DbSSLMode)
	l.Debug("数据库连接信息： ", dsn)
	_db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(io.Discard, "", 0),
			logger.Config{
				LogLevel: logger.Silent,
			},
		),
	})
	if err != nil {
		l.Fatal("数据库连接错误： ", err)
	}
	_db.AutoMigrate(&Sms{})
	_db.AutoMigrate(&Device{})
	_db.AutoMigrate(&User{})
	db = _db
}
