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

type KeyCloakConfig struct {
	BaseURL        string
	Realm          string
	ClientID       string `mapstructure:"client_id"`
	ClientSecret   string `mapstructure:"client_secret"`
	RealmPublicKey string `mapstructure:"realm_public_key"`
	Frontend       KeyCloakFrontendConfig
}

type KeyCloakFrontendConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	AuthURL      string `mapstructure:"auth_url"`
	TokenURL     string `mapstructure:"token_url"`
}

type IdentityAuthConfig struct {
	KeyCloak KeyCloakConfig
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
