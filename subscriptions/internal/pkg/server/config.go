package server

import "time"

type Config struct {
	Address           string        `yaml:"address" env-default:"localhost:8080"`
	Timeout           time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout       time.Duration `yaml:"idleTimeout" env-default:"60s"`
	ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout" env-default:"10s"`
}
