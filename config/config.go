package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		PG       pg
		Log      log
		Telegram telegram
	}

	pg struct {
		PoolMax  int    `env:"PG_POOL_MAX,required"`
		Host     string `env:"DB_HOST,required"`
		Port     string `env:"DB_PORT,required"`
		User     string `env:"DB_USER,required"`
		Password string `env:"DB_PASSWORD,required"`
		Name     string `env:"DB_NAME,required"`
		SSLMode  string `env:"DB_SSL_MODE" env-default:"disable"`
	}

	log struct {
		Level string `env:"LOG_LEVEL" env-default:"info"`
	}

	telegram struct {
		Token         string        `env:"TELEGRAM_TOKEN" env-default:""`
		PollerLimit   int           `env:"TELEGRAM_POLLER_LIMIT" env-default:"100"`
		PollerTimeout time.Duration `env:"TELEGRAM_POLLER_TIMEOUT" env-default:"10s"`
	}
)

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}

func (p *pg) URL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.User, p.Password, p.Host, p.Port, p.Name, p.SSLMode)
}
