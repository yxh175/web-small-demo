package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

// NewRedisCache 创建一个新的 Redis 缓存实例
func NewRedisCache() *RedisCache {
	return &RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379", // Redis 服务器地址
			Password: "",               // 如果有密码，设置密码
			DB:       0,                // 默认数据库
			PoolSize: 1000,
		}),
	}
}

// Ping 用于测试 Redis 连接
func (rc *RedisCache) Ping() error {
	ctx := context.Background()
	pong, err := rc.client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("无法连接到 Redis: %v", err)
	}
	fmt.Println("Redis 连接成功:", pong)
	return nil
}

// Set 用于设置缓存数据
func (rc *RedisCache) Set(key, value string, expiration time.Duration) error {
	ctx := context.Background()
	err := rc.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("设置缓存失败: %v", err)
	}
	return nil
}

// Get 用于获取缓存数据
func (rc *RedisCache) Get(key string) (string, error) {
	ctx := context.Background()
	cacheValue, err := rc.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", err
	} else if err != nil {
		return "", fmt.Errorf("获取缓存失败: %v", err)
	}
	return cacheValue, nil
}

// Delete 用于删除缓存数据
func (rc *RedisCache) Delete(key string) error {
	ctx := context.Background()
	err := rc.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("删除缓存失败: %v", err)
	}
	return nil
}
