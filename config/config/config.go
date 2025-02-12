package config

import (
	"log"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/Alias1177/merch-store/pkg"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
}

type AppConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}

func Load(path string) Config {
	// Читаем файл
	file, err := os.ReadFile(path)
	if err != nil {
		slog.Error(pkg.CfgErr)
		log.Fatalf("Unable to read config: %v", err)
	}

	// Парсим YAML
	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		slog.Error(pkg.ParseCfgErr)
		log.Fatalf("Failed to parse config: %v", err)
	}

	return cfg
}
