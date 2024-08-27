package config

import (
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

func InitViper() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
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
