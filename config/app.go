package config

import (
	"log/slog"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
)

type DatabaseConfig struct {
	Host     string        `mapstructure:"host"`
	Port     int           `mapstructure:"port"`
	Username string        `mapstructure:"username"`
	Password string        `mapstructure:"password"`
	Name     string        `mapstructure:"name"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

type MQTTConfig struct {
	Broker   string `mapstructure:"broker"`
	ClientID string `mapstructure:"client_id"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Topic    string `mapstructure:"topic"`
}

type LogConfig struct {
	Level  logger.LogLevel  `mapstructure:"level"`
	Format logger.LogFormat `mapstructure:"format"`
}

type ServerConfig struct {
	Logs        LogConfig      `mapstructure:"logs"`
	Database    DatabaseConfig `mapstructure:"database"`
	Port        int            `mapstructure:"port"`
	Development bool           `mapstructure:"development"`
	AppURL      string         `mapstructure:"app_url"`
}

type DashboardConfig struct {
	Title string `mapstructure:"title"`
}

type IdentityAuthConfig struct {
	OidcProvider OidcProvider `mapstructure:"oidc_provider"`
}

type Config struct {
	Server       ServerConfig       `mapstructure:"server"`
	Dashboard    DashboardConfig    `mapstructure:"dashboard"`
	MQTT         MQTTConfig         `mapstructure:"mqtt"`
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
