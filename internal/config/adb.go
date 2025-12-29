package config

import adb "github.com/zach-klippenstein/goadb"

var Adb = AdbConfig{}

type AdbConfig struct {
	AdbHost string
	AdbPort int
}

func (c *AdbConfig) initConfig() {
	c.AdbHost = getEnvOrDefault("ADB_HOST", "localhost")
	c.AdbPort = getIntEnv("ADB_PORT", adb.AdbPort)
}
