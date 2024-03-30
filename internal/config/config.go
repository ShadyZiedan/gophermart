package config

import (
	"flag"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func ParseConfig() (*Config, error) {
	var cfg Config
	flag.StringVar(&cfg.RunAddress, "a", ":8080", "The address to run on")
	flag.StringVar(&cfg.DatabaseURI, "d", "", "The database URI")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", "", "The accrual system address to run on")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
