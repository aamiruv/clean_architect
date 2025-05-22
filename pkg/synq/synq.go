// Package synq provides a seamless way to synchronize a cache with a database
// in a database-first manner. It ensures consistency by always prioritizing the database
// as the source of truth, while keeping the cache updated for performance.
//
// Key Features:
//   - **Database-First**: All operations (Get/Set/Delete) first modify the database
//     (via provided callbacks) before updating the cache.
//   - **Sync & Async Modes**: Supports both blocking (strong consistency) and
//     non-blocking (eventual consistency) operations.
//   - **CRUD Integration**: Wraps database calls with automatic cache management,
//     reducing boilerplate for read-through/write-through patterns.
//
// Usage Example:
//
//	store := synq.New[User](redisClient)
//	user, err := store.Get(ctx, "user123", func() (User, error) {
//	    return db.GetUserByID(ctx, "user123") // Database fetch
//	})
//
// Use Async methods (e.g., SetAsync) for write-heavy workloads where
// latency matters more than immediate consistency.
package synq

import (
	"context"
	"log/slog"

	"github.com/amirzayi/clean_architect/pkg/cache"
)

// CacheSync provides synchronized CRUD operations with caching support.
// It ensures cache consistency by invoking the given functions (getFn, setFn, deleteFn)
// to sync data between the cache and the underlying storage (e.g., DB, API).
type CacheSync[T any] interface {
	// Get retrieves a value by key. If the key is not in the cache (or expired),
	// it calls getFn to fetch the data, stores it in the cache, and returns it.
	//
	// Args:
	//   - ctx: Context for cancellation/timeout.
	//   - key: Cache key to retrieve.
	//   - getFn: Function that fetches data if cache misses.
	//
	// Returns:
	//   - The cached or freshly fetched value.
	//   - Error if the operation fails (e.g., getFn fails or cache is unreachable).
	Get(ctx context.Context, key string, getFn func() (T, error)) (T, error)

	// Set updates a key in the cache and optionally syncs it to the underlying storage via setFn.
	// If setFn is provided, it is called to persist the data before updating the cache.
	//
	// Args:
	//   - ctx: Context for cancellation/timeout.
	//   - key: Cache key to update.
	//   - value: Value to store.
	//   - setFn: function to persist the data (e.g., DB write).
	//
	// Returns:
	//   - Error if the operation fails (e.g., setFn fails or cache write fails).
	Set(ctx context.Context, key string, value T, setFn func() error) error

	// Delete removes a key from the cache for cache invalidation after calling deleteFn.
	//
	// Args:
	//   - ctx: Context for cancellation/timeout.
	//   - key: Cache key to delete.
	//   - deleteFn: function to delete data from storage (e.g., DB delete).
	//
	// Returns:
	//   - Error if the operation fails (e.g., deleteFn fails or cache deletion fails).
	Delete(ctx context.Context, key string, deleteFn func() error) error

	// GetAsync is similar to Get but executes cache fulfillment asynchronously.
	//
	// Args:
	//   - ctx: Context for cancellation/timeout.
	//   - key: Cache key to retrieve.
	//   - getFn: Function that fetches data if cache misses (executed in background).
	//
	// Returns:
	//   - The cached or freshly fetched value.
	//   - Error only if the getFn fails (cache fetch errors are logged).
	GetAsync(ctx context.Context, key string, getFn func() (T, error)) (T, error)

	// SetAsync is similar to Set but executes cache fulfillment asynchronously.
	//
	// Args:
	//   - key: Cache key to update.
	//   - value: Value to store.
	//   - setFn: function to persist the data (e.g., DB write).
	//
	// Returns:
	//   - Error only if the cache update fails (cache fulfillment errors are logged).
	SetAsync(key string, value T, setFn func() error) error

	// DeleteAsync is similar to Delete but executes cache removal asynchronously.
	//
	// Args:
	//   - key: Cache key to delete.
	//   - deleteFn: Function to delete data from storage.
	//
	// Returns:
	//   - Error only if the cache deletion fails (cache removal errors are logged).
	DeleteAsync(key string, deleteFn func() error) error
}

type synq[T any] struct {
	cache  cache.Cache[T]
	logger *slog.Logger
}

func New[T any](cache cache.Cache[T], logger *slog.Logger) CacheSync[T] {
	return synq[T]{
		cache:  cache,
		logger: logger,
	}
}
