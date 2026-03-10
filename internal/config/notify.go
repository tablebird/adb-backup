package config

var Notify = NotifyConfig{}

type NotifyConfig struct {
	NotifyWebhookUrl       string
	NotifyStatusWebhookUrl string
}

func (c *NotifyConfig) initConfig() {
	c.NotifyWebhookUrl = getEnvOrDefault("NOTIFY_WEBHOOK_URL", "")
	c.NotifyStatusWebhookUrl = getEnvOrDefault("NOTIFY_STATUS_WEBHOOK_URL", "")
}
