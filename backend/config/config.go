package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
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

// Init 读取配置。
//
// 有两种启动方式：
//   - 手动启动（go run . / 本地编译运行）：默认会加载当前目录的 .env（由 .env.example 复制而来），
//     文件不存在时直接使用环境变量与下面的默认值。
//   - 容器启动：镜像 Dockerfile 里设置了 LOAD_DOTENV=false，因此完全不读取 .env，
//     配置只来自镜像 ENV（docker-compose.yml 不再重复声明）。
//
// 可用 LOAD_DOTENV=false 显式关闭 .env 加载；无论哪种方式，未设置的项都使用默认值。
func Init() error {
	if getEnv("LOAD_DOTENV", "true") != "false" {
		if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	redisDB, err := getEnvInt("REDIS_DB", 0)
	if err != nil {
		return err
	}

	AppConfig = Config{
		RedisAddr:     fmt.Sprintf("%s:%s", getEnv("REDIS_HOST", "localhost"), getEnv("REDIS_PORT", "6379")),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,
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

// getEnvInt 解析整数环境变量；值非法时返回错误，而不是静默回退默认值
func getEnvInt(key string, defaultValue int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s=%q: %w", key, value, err)
	}
	return parsed, nil
}
