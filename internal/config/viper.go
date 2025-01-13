package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func InitViper() (*Config, error) {
	osEnv := os.Getenv("ENV")
	if osEnv == "" {
		osEnv = "dev"
	}
	viper.SetConfigName(fmt.Sprintf("config.%s", osEnv))
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.SetEnvPrefix("GE")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		slog.Error("Error unmarshalling config", "error", err)
		return nil, err
	}

	return &cfg, nil
}
