package config

import (
	"fmt"

	"TZshka/pkg/postgres"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port     string          `yaml:"PORT"`
	Secret   string          `yaml:"JWT_SECRET"`
	LLMURL   string          `yaml:"LLM_API_URL"`
	Postgres postgres.Config `yaml:"POSTGRES"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig("config/config.yaml", &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	return &cfg, nil
}
