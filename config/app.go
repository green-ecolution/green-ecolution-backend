package config

import (
	"log/slog"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Name     string
	Timeout  time.Duration
}

type MQTTConfig struct {
	Broker   string
	ClientID string `mapstructure:"client_id"`
	Username string
	Password string
	Topic    string
}

type LogConfig struct {
	Level  logger.LogLevel
	Format logger.LogFormat
}

type ServerConfig struct {
	Logs        LogConfig
	Database    DatabaseConfig
	Port        int
	Development bool
	AppURL      string `mapstructure:"app_url"`
}

type DashboardConfig struct {
	Title string
}

type IdentityAuthConfig struct {
  OidcProvider OidcProvider `mapstructure:"oidc_provider"`
}

type Config struct {
	Server       ServerConfig
	Dashboard    DashboardConfig
	MQTT         MQTTConfig
	IdentityAuth IdentityAuthConfig `mapstructure:"auth"`
}

func InitConfig() (*Config, error) {
	slog.Info("Loading config...")

	cfg, err := InitViper()
	if err != nil {
		return nil, err
	}

	slog.Info("Config loaded successfully")
	return cfg, nil
}
