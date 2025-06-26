package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	App      AppConfig
}

type DatabaseConfig struct {
	URL string
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type JWTConfig struct {
	Secret string
}

type AppConfig struct {
	Name    string
	Version string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/hospital_management?sslmode=disable"),
		},
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "default_secret"),
		},
		App: AppConfig{
			Name:    getEnv("APP_NAME", "Hospital Management System"),
			Version: getEnv("APP_VERSION", "1.0.0"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
