package config

import (
	"log"
	"log/slog"

	"github.com/Alias1177/merch-store/pkg"
	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	Port string `env:"APP_PORT" env-default:"8080"`
}

type DatabaseConfig struct {
	DSN string `env:"DATABASE_DSN" env-required:"true"`
}

type JWTConfig struct {
	Secret string `env:"JWT_SECRET" env-required:"true"`
}

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

func Load(path string) Config {
	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		slog.Error(pkg.CfgErr)
		log.Fatalf("Unable to read config: %v", err)
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		slog.Error(pkg.CfgErr)
		log.Fatalf("Unable to read environment variables: %v", err)
	}

	return cfg
}
