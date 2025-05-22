package synq

import (
	"context"
	"log/slog"
)

func (cr synq[T]) GetAsync(ctx context.Context, key string, getFn func() (T, error)) (T, error) {
	data, err := cr.cache.Get(ctx, key)
	if err == nil {
		return data, nil
	}
	data, err = getFn()
	if err != nil {
		return data, err
	}
	go func() {
		if err = cr.cache.Set(context.Background(), key, data); err != nil {
			cr.logger.Error("failed to set cache", slog.String("key", key), slog.Any("error", err))
		}
	}()
	return data, nil
}

func (cr synq[T]) SetAsync(key string, value T, setFn func() error) error {
	if err := setFn(); err != nil {
		return err
	}
	go func() {
		if err := cr.cache.Set(context.Background(), key, value); err != nil {
			cr.logger.Error("failed to set cache", slog.String("key", key), slog.Any("error", err))
		}
	}()
	return nil
}

func (cr synq[T]) DeleteAsync(key string, deleteFn func() error) error {
	if err := deleteFn(); err != nil {
		return err
	}
	go func() {
		if err := cr.cache.Delete(context.Background(), key); err != nil {
			cr.logger.Error("failed to delete cache", slog.String("key", key), slog.Any("error", err))
		}
	}()
	return nil
}
