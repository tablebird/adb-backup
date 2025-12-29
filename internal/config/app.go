package config

import (
	"time"
)

var App = AppConfig{}

type AppConfig struct {
	DebugLog bool
	// 同步等待间隔
	ReadInterval       time.Duration
	WaitDeviceInterval time.Duration

	NotifyWebhookUrl string
}

func (c *AppConfig) initConfig() {
	c.DebugLog = getBoolEnv("DEBUG_LOG", false)
	c.ReadInterval = time.Second * time.Duration(getIntEnv("READ_INTERVAL", 5))
	c.WaitDeviceInterval = time.Second * time.Duration(getIntEnv("WAIT_DEVICE_INTERVAL", 10))
	c.NotifyWebhookUrl = getEnvOrDefault("NOTIFY_WEBHOOK_URL", "")
}
