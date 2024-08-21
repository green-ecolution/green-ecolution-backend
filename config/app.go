package config

import (
	"errors"
	"log/slog"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string        `env:"HOST" envDefault:"localhost"`
	Port     int           `env:"PORT" envDefault:"27017"`
	Username string        `env:"USER" envDefault:"root"`
	Password string        `env:"PASSWORD" envDefault:"example"`
	Name     string        `env:"NAME" envDefault:"green-space-management"`
	Timeout  time.Duration `env:"TIMEOUT" envDefault:"10s"`
}

type MQTTConfig struct {
	Broker   string `env:"BROKER" envDefault:"eu1.cloud.thethings.network:1883"`
	ClientID string `mapstructure:"client_id" env:"CLIENT_ID"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	Topic    string `env:"TOPIC"`
}

type LogConfig struct {
	Level  logger.LogLevel
	Format logger.LogFormat
}

type ServerConfig struct {
	Logs        LogConfig      `envPrefix:"LOGS"`
	Database    DatabaseConfig `envPrefix:"DATABASE"`
	Port        int            `env:"PORT" envDefault:"3000"`
	Development bool           `env:"DEVELOPMENT" envDefault:"false"`
	AppURL      string         `mapstructure:"app_url" env:"APP_URL" envDefault:"http://localhost:$PORT"`
}

type DashboardConfig struct {
	Title string `env:"TITLE" envDefault:"Green Ecolution Dashboard"`
}

type KeyCloakConfig struct {
	BaseURL        string `env:"BASE_URL"`
	Realm          string `env:"REALM"`
	ClientID       string `mapstructure:"client_id" env:"CLIENT_ID"`
	ClientSecret   string `mapstructure:"client_secret" env:"CLIENT_SECRET"`
	RealmPublicKey string `mapstructure:"realm_public_key" env:"REALM_PUBLIC_KEY"`
}

type IdentityAuthConfig struct {
	KeyCloak KeyCloakConfig `envPrefix:"KEYCLOAK_"`
}

type Config struct {
	Server       ServerConfig       `envPrefix:"GE_SERVER_"`
	Dashboard    DashboardConfig    `envPrefix:"GE_DASHBOARD_"`
	MQTT         MQTTConfig         `envPrefix:"GE_MQTT_"`
	IdentityAuth IdentityAuthConfig `mapstructure:"auth" envPrefix:"GE_AUTH_"`
}

var (
	ErrViperConfigFileNotFound = viper.ConfigFileNotFoundError{}
	ErrViperConfigFileError    = errors.New("error loading config file with viper")
	ErrEnvConfigError          = errors.New("error loading config from environment variables")
)

func InitConfig() (*Config, error) {
	slog.Info("Loading config...")

	cfg, err := InitViper()
	if err != nil {
		if errors.Is(err, ErrViperConfigFileNotFound) {
			slog.Info("Config file not found, trying to load from environment variables")
			cfg, err = InitEnv()
			if err != nil {
				slog.Error("Error loading config from environment variables", "error", err)
				return nil, errors.Join(err, ErrEnvConfigError)
			}
		} else {
			slog.Error("Error loading config file", "error", err)
			return nil, err
		}
	}

	slog.Info("Config loaded successfully")
	return cfg, nil
}
