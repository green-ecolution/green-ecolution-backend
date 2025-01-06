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

type RoutingConfig struct {
	StartPoint    []float64        `mapstructure:"start_point"`
	EndPoint      []float64        `mapstructure:"end_point"`
	WateringPoint []float64        `mapstructure:"watering_point"`
	Ors           RoutingOrsConfig `mapstructure:"ors"`
}

type RoutingOrsConfig struct {
	Host         string                       `mapstructure:"host"`
	Optimization RoutingOrsOptimizationConfig `mapstructure:"optimization"`
}

type RoutingOrsOptimizationConfig struct {
	Vroom RoutingVroomConfig `mapstructure:"vroom"`
}

type RoutingVroomConfig struct {
	Host string `mapstructure:"host"`
}

type S3Config struct {
	Endpoint string          `mapstructure:"endpoint"`
	Region   string          `mapstructure:"region"`
	RouteGpx S3ServiceConfig `mapstructure:"route-gpx"`
}

type S3ServiceConfig struct {
	Bucket          string `mapstructure:"bucket"`
	AccessKey       string `mapstructure:"accessKey"`
	SecretAccessKey string `mapstructure:"secretAccessKey"`
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
	Routing      RoutingConfig      `mapstructure:"routing"`
	S3           S3Config           `mapstructure:"s3"`
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
