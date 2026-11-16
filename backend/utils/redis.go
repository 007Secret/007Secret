package utils

import (
	"context"
	"encoding/base64"

	"github.com/007Secret/007Secret/config"
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.RedisAddr,
		Password: config.AppConfig.RedisPassword,
		DB:       config.AppConfig.RedisDB,
	})

	_, err := RedisClient.Ping(ctx).Result()
	return err
}

func StoreSecret(key string, value []byte) error {
	// 使用 base64 编码压缩后的数据
	encodedValue := base64.StdEncoding.EncodeToString(value)
	return RedisClient.Set(ctx, key, encodedValue, config.AppConfig.ExpireTime).Err()
}

func GetSecret(key string) ([]byte, error) {
	// 获取并删除数据
	value, err := RedisClient.GetDel(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	// 解码数据
	return base64.StdEncoding.DecodeString(value)
}
