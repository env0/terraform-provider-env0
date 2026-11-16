package client

import "sync"

type memoizedResult[V any] struct {
	once  sync.Once
	value V
	err   error
}

// Terraform reads data sources in parallel, so the cache is locked. The lock covers only the map:
// each key's Once makes concurrent callers for that key share one call to f.
func memoize[K comparable, V any](f func(K) (V, error)) func(K) (V, error) {
	var mu sync.Mutex

	cache := make(map[K]*memoizedResult[V])

	return func(key K) (V, error) {
		mu.Lock()

		res, ok := cache[key]
		if !ok {
			res = &memoizedResult[V]{}
			cache[key] = res
		}

		mu.Unlock()

		res.once.Do(func() {
			res.value, res.err = f(key)
		})

		return res.value, res.err
	}
}

// MemoizeExported exports the memoize function for testing
func MemoizeExported[K comparable, V any](f func(K) (V, error)) func(K) (V, error) {
	return memoize(f)
}
