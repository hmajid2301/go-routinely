package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DB     DBConfig
	Server ServerConfig
}

type DBConfig struct {
	URI string
}

type ServerConfig struct {
	Host string
	Port int
}

func LoadConfig(ctx context.Context) (*Config, error) {
	dbURI := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/goroutinely?sslmode=disable")
	serverHost := getEnv("GOROUTINELY_WEBSERVER_HOST", "0.0.0.0")
	serverPort, err := strconv.Atoi(getEnv("GOROUTINELY_WEBSERVER_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	return &Config{
		DB: DBConfig{
			URI: dbURI,
		},
		Server: ServerConfig{
			Host: serverHost,
			Port: serverPort,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
