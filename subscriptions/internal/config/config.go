package config

import (
	"github.com/ekkserapopova/subscriptions/internal/pkg/db"
	"github.com/ekkserapopova/subscriptions/internal/pkg/server"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"log"
	"os"
)

type Config struct {
	ConfigPath string `env:"CONFIG_PATH" env-default:"config/config.yaml"`

	HTTPServer server.Config `yaml:"httpServer"`
	DB         db.Config     `yaml:"db"`
}

type Out struct {
	fx.Out

	HTTPServer server.Config
	DB         db.Config
}

func MustLoad() Out {
	// 1. Подгружаем переменные из .env, если файл существует
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("failed to load .env file: %s", err)
			os.Exit(1)
		}
	}

	var cfg Config

	// 2. Читаем CONFIG_PATH из env
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Printf("cannot read env variables: %s", err)
		os.Exit(1)
	}

	// 3. Читаем YAML
	if _, err := os.Stat(cfg.ConfigPath); os.IsNotExist(err) {
		log.Printf("config file does not exist: %s", cfg.ConfigPath)
		os.Exit(1)
	}
	if err := cleanenv.ReadConfig(cfg.ConfigPath, &cfg); err != nil {
		log.Printf("cannot read %s: %v", cfg.ConfigPath, err)
		os.Exit(1)
	}

	// 4. Переопределяем настройками из env (POSTGRES_*)
	if err := cleanenv.ReadEnv(&cfg.DB); err != nil {
		log.Printf("cannot read DB env variables: %s", err)
		os.Exit(1)
	}
	if err := cleanenv.ReadEnv(&cfg.HTTPServer); err != nil {
		log.Printf("cannot read HTTPServer env variables: %s", err)
		os.Exit(1)
	}

	return Out{
		HTTPServer: cfg.HTTPServer,
		DB:         cfg.DB,
	}
}
