//go:build dev

package config

import (
	"flag"
	"log"

	godotenv "github.com/joho/godotenv"
)

func initEnv() {
	// 加载 .env 文件
	var err error
	if flag.Lookup("test.v") == nil {
		err = godotenv.Load(".env")
	} else {
		err = godotenv.Load("../../.env")
	}
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
