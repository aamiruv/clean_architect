package cache

import (
	"context"
	"sync"
	"time"
)

type inMemCache struct {
	mu    sync.RWMutex
	store map[string]item
}

type item struct {
	data      []byte
	expiresAt time.Time
}

func NewInMemoryDriver() Driver {
	return &inMemCache{store: make(map[string]item)}
}

func (c *inMemCache) Set(_ context.Context, key string, data []byte, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	exp := time.Time{}
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}
	c.store[key] = item{data: data, expiresAt: exp}
	return nil
}

func (c *inMemCache) Get(ctx context.Context, key string) ([]byte, error) {
	c.mu.RLock()
	itm, ok := c.store[key]
	c.mu.RUnlock()

	if !ok {
		return nil, ErrCacheMissed
	}

	if !itm.expiresAt.IsZero() && time.Now().After(itm.expiresAt) {
		c.Delete(ctx, key)
		return nil, ErrCacheMissed
	}
	return itm.data, nil
}

func (c *inMemCache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	delete(c.store, key)
	c.mu.Unlock()
	return nil
}

func (c *inMemCache) Close() error {
	return nil
}
