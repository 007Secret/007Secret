package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	ExpireTime    time.Duration
	Port          string
}

var AppConfig Config

func Init() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")

	AppConfig = Config{
		RedisAddr:     fmt.Sprintf("%s:%s", host, port),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       0,
		ExpireTime:    24 * time.Hour,
		Port:          getEnv("PORT", "8080"),
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
