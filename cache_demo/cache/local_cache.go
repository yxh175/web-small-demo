package cache

import (
	"sync"
	"time"
)

var LocalCache *Cache

// 缓存数据的结构
type Cache struct {
	mu      sync.RWMutex
	data    map[string]interface{}
	expires map[string]time.Time
}

// 创建一个新的缓存
func NewCache() *Cache {
	return &Cache{
		data:    make(map[string]interface{}),
		expires: make(map[string]time.Time),
	}
}

// 向缓存中添加数据
func (c *Cache) Set(key string, value interface{}, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	c.expires[key] = time.Now().Add(expiration)
}

// 从缓存中获取数据
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.data[key]
	if !ok {
		return nil, false
	}
	expiration, exists := c.expires[key]
	if !exists || time.Now().Before(expiration) {
		return value, true
	}
	// 数据已过期，从缓存中删除
	delete(c.data, key)
	delete(c.expires, key)
	return nil, false
}

func InitLocal() {
	LocalCache = NewCache()
}
