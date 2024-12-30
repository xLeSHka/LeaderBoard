package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	cache1 "github.com/xLeSHka/LeaderBoard/pkg/db/cache"
)

type Config struct {
	ServerPort int `env:"SERVER_PORT" env-default:"9090"`
	cache1.RedisConfig
}

func New() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
