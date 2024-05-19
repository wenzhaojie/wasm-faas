package localcache

import (
	"sync"
)

// Cache 是一个简单的本地缓存对象
type Cache struct {
	data map[string]interface{}
	mu   sync.Mutex
}

// NewCache 创建一个新的缓存对象
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
	}
}

// Put 将键值对放入缓存中
func (c *Cache) Put(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

// Get 从缓存中获取键对应的值
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.data[key]
	return value, ok
}
