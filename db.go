package main

import (
	"fmt"
	"io"
	"log"

	postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.DbHost, config.DbPort, config.DbUser, config.DbPass, config.DbName, config.DbSSLMode)
	logDebug("数据库连接信息： ", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(io.Discard, "", 0),
			logger.Config{
				LogLevel: logger.Silent,
			},
		),
	})
	if err != nil {
		logFatal("数据库连接错误： ", err)
	}
	db.AutoMigrate(&Sms{})
	db.AutoMigrate(&Device{})
	return db
}
