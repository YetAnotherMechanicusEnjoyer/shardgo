package cache

import (
	"errors"
	"sync"
	"time"
)

var ErrKeyNotFound = errors.New("key not found")

type CacheEntry struct {
	Value     []byte
	ExpiresAt time.Time
}

type Cache struct {
	data sync.Map
}

func New() *Cache {
	return &Cache{}
}

func (c *Cache) Set(key string, value []byte, ttl time.Duration) error {
	entry := &CacheEntry{
		Value: value,
	}
	if ttl > 0 {
		entry.ExpiresAt = time.Now().Add(ttl)
	}
	c.data.Store(key, entry)
	return nil
}

func (c *Cache) Get(key string) ([]byte, error) {
	entry, ok := c.data.Load(key)
	if !ok {
		return nil, ErrKeyNotFound
	}
	e := entry.(*CacheEntry)
	if !e.ExpiresAt.IsZero() && time.Now().After(e.ExpiresAt) {
		c.data.Delete(key)
		return nil, ErrKeyNotFound
	}
	return e.Value, nil
}

func (c *Cache) Delete(key string) error {
	_, ok := c.data.Load(key)
	if !ok {
		return ErrKeyNotFound
	}
	c.data.Delete(key)
	return nil
}

func (c *Cache) Size() int {
	count := 0
	c.data.Range(func(key, value any) bool {
		count++
		return true
	})
	return count
}
