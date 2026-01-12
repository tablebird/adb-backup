package config

var Feature = FeatureConfig{}

type FeatureConfig struct {
	EnableSendSms bool
}

func (c *FeatureConfig) initConfig() {
	c.EnableSendSms = getBoolEnv("ENABLE_SEND_SMS", true)
}
