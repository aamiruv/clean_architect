package synq

import "context"

func (cr synq[T]) Get(ctx context.Context, key string, getFn func() (T, error)) (T, error) {
	data, err := cr.cache.Get(ctx, key)
	if err == nil {
		return data, nil
	}
	data, err = getFn()
	if err == nil {
		cr.cache.Set(ctx, key, data)
	}
	return data, err
}

func (cr synq[T]) Set(ctx context.Context, key string, value T, setFn func() error) error {
	if err := setFn(); err != nil {
		return err
	}
	return cr.cache.Set(ctx, key, value)
}

func (cr synq[T]) Delete(ctx context.Context, key string, deleteFn func() error) error {
	if err := deleteFn(); err != nil {
		return err
	}
	return cr.cache.Delete(ctx, key)
}
