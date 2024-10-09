package cache

import (
	"sync"
)

type InMemoryCache struct {
	data map[string]any
	lock sync.RWMutex
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]any),
	}
}

func (c *InMemoryCache) Set(key string, value any) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data[key] = value
	return nil
}

func (c *InMemoryCache) Get(key string) (any, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if value, ok := c.data[key]; ok {
		return value, nil
	} else {
		return nil, nil
	}
}

func (c *InMemoryCache) Del(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.data, key)
	return nil
}
