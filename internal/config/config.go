package config

type Config struct {
	Service  *ServiceConfig  `mapstructure:"service"`
	Database *DatabaseConfig `mapstructure:"database"`
}
