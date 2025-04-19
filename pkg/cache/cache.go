package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"time"
)

var (
	ErrCacheMissed = errors.New("cache missed")
)

type Driver interface {
	Set(ctx context.Context, key string, data []byte, ttl time.Duration) error
	Get(ctx context.Context, key string) (data []byte, err error)
	Delete(ctx context.Context, key string) error
}

type Cache[T any] interface {
	Set(ctx context.Context, key string, value T, ttl time.Duration) error
	Get(ctx context.Context, key string) (value T, err error)
	Delete(ctx context.Context, key string) error
}

type typedCache[T any] struct {
	drv Driver
}

func New[T any](drv Driver) Cache[T] {
	return &typedCache[T]{drv: drv}
}

func (c *typedCache[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return err
	}
	return c.drv.Set(ctx, key, buf.Bytes(), ttl)
}

func (c *typedCache[T]) Get(ctx context.Context, key string) (T, error) {
	var v T

	b, err := c.drv.Get(ctx, key)
	if err != nil {
		return v, err
	}

	buf := bytes.NewBuffer(b)
	if err = gob.NewDecoder(buf).Decode(&v); err != nil {
		return v, err
	}
	return v, nil
}

func (c *typedCache[T]) Delete(ctx context.Context, key string) error {
	return c.drv.Delete(ctx, key)
}
