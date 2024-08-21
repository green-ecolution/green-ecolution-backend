package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func InitEnv() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
