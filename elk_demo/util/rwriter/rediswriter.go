package rwriter

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

var rWriter *RedisWriter
var once sync.Once

type RedisWriter struct {
	cli     *redis.Client
	listKey string
	c       context.Context
}

func (w *RedisWriter) Write(p []byte) (int, error) {
	n, err := w.cli.RPush(w.c, w.listKey, p).Result()
	return int(n), err
}

// NewRedisWriter
func NewRedisWriter() *RedisWriter {
	once.Do(func() {
		rWriter = &RedisWriter{
			cli: redis.NewClient(&redis.Options{
				Addr:     "localhost:6379", // Redis 服务器地址
				Password: "",               // 如果有密码，设置密码
				DB:       0,                // 默认数据库
				PoolSize: 1000,
			}),
			listKey: "log_queue",
			c:       context.Background(),
		}
	})
	return rWriter
}
