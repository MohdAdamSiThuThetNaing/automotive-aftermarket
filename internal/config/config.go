package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL string
	ProductID   string
}

func Load() *Config {

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		ProductID:   os.Getenv("PRODUCT_ID"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	return cfg
}