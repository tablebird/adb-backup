//go:build dev

package config

import (
	"log"

	godotenv "github.com/joho/godotenv"
)

func initEnv() {
	// 加载 .env 文件
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
