package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisURL       string
	PostgresURL    string
	HTTPServerAddr string
}

func LoadEnv() {
	if godotenv.Load() != nil {
		log.Fatal("Error loading .env file")
	}
}

func Load() *Config {
	LoadEnv()
	return &Config{
		PostgresURL:    GetEnv("POSTGRES_URL"),
		RedisURL:       GetEnv("REDIS_URL"),
		HTTPServerAddr: GetEnv("HTTP_SERVER_ADDRESS"),
	}
}

func GetEnv(key string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Missing required environment variable: %s. Please set it in your environment or .env file.", key)
	}
	return val
}
