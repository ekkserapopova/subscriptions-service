package db

import "time"

type Config struct {
	DB             string        `env:"POSTGRES_DB"`
	User           string        `env:"POSTGRES_USER"`
	Password       string        `env:"POSTGRES_PASSWORD"`
	Host           string        `env:"POSTGRES_HOST"`
	Port           uint16        `env:"POSTGRES_PORT"`
	ConnectTimeout time.Duration `yaml:"connectTimeout" env-default:"5m"`
}
