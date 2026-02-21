package memory

import (
	"context"
	"sync"
)

type Cache struct {
	mu    sync.RWMutex
	store map[string]string
}

func New() *Cache {
	return &Cache{store: map[string]string{}}
}

func (c *Cache) Get(_ context.Context, key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.store[key]
	return v, ok
}

func (c *Cache) Set(_ context.Context, key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
}

func (c *Cache) Delete(_ context.Context, key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

func (c *Cache) Exists(ctx context.Context, key string) bool {
	_, ok := c.Get(ctx, key)
	return ok
}
