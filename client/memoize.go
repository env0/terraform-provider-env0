package client

type memoizedResult[V any] struct {
	value V
	err   error
}

func memoize[K comparable, V any](f func(K) (V, error)) func(K) (V, error) {
	cache := make(map[K]memoizedResult[V])

	return func(key K) (V, error) {
		if res, ok := cache[key]; ok {
			return res.value, res.err
		}

		value, err := f(key)

		cache[key] = memoizedResult[V]{value: value, err: err}

		return value, err
	}
}
