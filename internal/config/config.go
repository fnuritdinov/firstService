package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HttpPort string `env:"HTTP_PORT"`
	Storage  string `env:"STORAGE"`
}

func New(filePath string) (Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(filePath, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("cleanenv.ReadConfig %w", err)
	}
	return cfg, err
}
