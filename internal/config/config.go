package config

import (
	"os"
	"strconv"
	"time"

	adb "github.com/zach-klippenstein/goadb"
)

var Conf = Config{}

type Config struct {
	AdbHost   string
	AdbPort   int
	DbHost    string
	DbPort    int
	DbUser    string
	DbPass    string
	DbName    string
	DbSSLMode string
	DebugLog  bool
	// 同步等待间隔
	ReadInterval       time.Duration
	WaitDeviceInterval time.Duration

	NotifyWebhookUrl string
}

func (c *Config) InitConfig() {
	initEnv()
	c.AdbHost = getEnvOrDefault("ADB_HOST", "localhost")
	c.AdbPort = getIntEnv("ADB_PORT", adb.AdbPort)
	c.DbHost = getEnvOrDefault("DB_HOST", "postgres.lan")
	c.DbPort = getIntEnv("DB_PORT", 5432)
	c.DbUser = getEnvOrDefault("DB_USER", "backup")
	c.DbPass = getEnvOrDefault("DB_PASS", "backup")
	c.DbName = getEnvOrDefault("DB_NAME", "backup")
	c.DbSSLMode = getEnvOrDefault("DB_SSLMODE", "disable")
	c.DebugLog = getBoolEnv("DEBUG_LOG", false)
	c.ReadInterval = time.Second * time.Duration(getIntEnv("READ_INTERVAL", 5))
	c.WaitDeviceInterval = time.Second * time.Duration(getIntEnv("WAIT_DEVICE_INTERVAL", 10))
	c.NotifyWebhookUrl = getEnvOrDefault("NOTIFY_WEBHOOK_URL", "")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
