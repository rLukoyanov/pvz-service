package config

import (
	"github.com/caarlos0/env"
)

type Config struct {
	DATABASE_URL string `env:"DATABASE_URL" envDefault:"postgresql://root:123@localhost:5432/base?sslmode=disable"`
	LOG_LEVEL    string `env:"LOG_LEVEL" envDefault:"debug"`
	MODE         string `env:"MODE" envDefault:"dev"`
	SECRET       string `env:"SECRET" envDefault:"your_secret_key"`
}

func NewConfig() (*Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
