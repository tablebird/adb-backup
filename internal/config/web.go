package config

var Web = WebConfig{}

type WebConfig struct {
	WebPort int

	AdminName     string
	AdminPassword string

	CookieMaxAge int
}

func (c *WebConfig) initConfig() {
	c.WebPort = getIntEnv("WEB_PORT", 8080)
	c.AdminName = getEnvOrDefault("ADMIN_NAME", "admin")
	c.AdminPassword = getEnvOrDefault("ADMIN_PASSWORD", "admin")
	c.CookieMaxAge = getIntEnv("COOKIE_MAX_AGE", 86400) //24h
}
