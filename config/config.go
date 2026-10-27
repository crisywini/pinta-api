package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBDriver    string
	MongoURI    string
	MongoDB     string
	PostgresURL string
}

func Load() (Config, error) {
	cfg := Config{
		DBDriver:    os.Getenv("DB"),
		MongoURI:    os.Getenv("MONGODB_URI"),
		MongoDB:     getenv("MONGODB_NAME", "pintapp"),
		PostgresURL: os.Getenv("DATABASE_URL"),
	}

	switch cfg.DBDriver {
	case "mongo":
		if cfg.MongoURI == "" {
			return cfg, fmt.Errorf("MONGODB_URI is required when DB=mongo")
		}
	case "postgres":
		if cfg.PostgresURL == "" {
			return cfg, fmt.Errorf("DATABASE_URL is required when DB=postgres")
		}
	case "":
		return cfg, fmt.Errorf("DB must be set (mongo|postgres)")
	}
	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
