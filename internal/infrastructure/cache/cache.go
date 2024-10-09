package cache

import "sync"

var (
	once  sync.Once
	cache Cache
)

type Cache interface {
	Set(key string, value any) error
	Get(key string) (any, error)
	Del(key string) error
}

func GetCache() Cache {
	once.Do(func() {
		cache = NewInMemoryCache()
	})
	return cache
}
