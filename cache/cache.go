package cache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string][]byte),
	}
}

func (c *Cache) Delete(key []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, string(key))
	return nil
}

func (c *Cache) Has(key []byte) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.data[string(key)]

	return ok
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, ok := c.data[string(key)]

	if !ok {
		return nil, fmt.Errorf("key (%s) not found", string(key))
	}

	return val, nil
}

func (c *Cache) Set(key, value []byte, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ttl > 0 {
		go func() {
			<-time.After(ttl)
			delete(c.data, string(key))
		}()
	}

	c.data[string(key)] = value

	return nil
}
