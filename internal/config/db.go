package config

var DB = DBConfig{}

type DBConfig struct {
	DbHost    string
	DbPort    int
	DbUser    string
	DbPass    string
	DbName    string
	DbSSLMode string
}

func (c *DBConfig) initConfig() {
	c.DbHost = getEnvOrDefault("DB_HOST", "postgres.lan")
	c.DbPort = getIntEnv("DB_PORT", 5432)
	c.DbUser = getEnvOrDefault("DB_USER", "backup")
	c.DbPass = getEnvOrDefault("DB_PASS", "backup")
	c.DbName = getEnvOrDefault("DB_NAME", "backup")
	c.DbSSLMode = getEnvOrDefault("DB_SSLMODE", "disable")
}
