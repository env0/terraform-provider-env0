package client_test

import (
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoize(t *testing.T) {
	t.Run("should cache successful results", func(t *testing.T) {
		callCount := 0
		f := func(s string) ([]client.Team, error) {
			callCount++

			return []client.Team{{Name: s}}, nil
		}

		memoized := client.MemoizeExported(f)

		// First call
		result1, err1 := memoized("")
		require.NoError(t, err1)
		assert.Len(t, result1, 1)
		assert.Equal(t, 1, callCount)

		// Second call with same input - should use cache
		result2, err2 := memoized("")
		require.NoError(t, err2)
		assert.Len(t, result2, 1)
		assert.Equal(t, 1, callCount) // Count shouldn't increase

		// Different input - should call function again
		result3, err3 := memoized("test")
		require.NoError(t, err3)
		assert.Len(t, result3, 1)
		assert.Equal(t, "test", result3[0].Name)
		assert.Equal(t, 2, callCount)
	})

	t.Run("should cache errors", func(t *testing.T) {
		callCount := 0
		expectedError := errors.New("test error")
		f := func(s string) ([]client.Team, error) {
			callCount++

			return nil, expectedError
		}

		memoized := client.MemoizeExported(f)

		// First call
		result1, err1 := memoized("")
		require.Error(t, err1)
		assert.Equal(t, expectedError, err1)
		assert.Nil(t, result1)
		assert.Equal(t, 1, callCount)

		// Second call - should return cached error
		result2, err2 := memoized("")
		require.Error(t, err2)
		assert.Equal(t, expectedError, err2)
		assert.Nil(t, result2)
		assert.Equal(t, 1, callCount) // Count shouldn't increase
	})

	t.Run("should be safe for concurrent calls", func(t *testing.T) {
		memoized := client.MemoizeExported(func(s string) ([]client.Team, error) {
			return []client.Team{{Name: s}}, nil
		})

		var wg sync.WaitGroup
		for i := range 1000 {
			name := strconv.Itoa(i % 100)
			wg.Go(func() {
				teams, err := memoized(name)
				assert.NoError(t, err)
				assert.Equal(t, []client.Team{{Name: name}}, teams)
			})
		}
		wg.Wait()
	})

	t.Run("should call the function once for concurrent calls with one key", func(t *testing.T) {
		var callCount atomic.Int32
		memoized := client.MemoizeExported(func(s string) ([]client.Team, error) {
			callCount.Add(1)
			time.Sleep(10 * time.Millisecond)

			return []client.Team{{Name: s}}, nil
		})

		var wg sync.WaitGroup
		for range 50 {
			wg.Go(func() {
				teams, err := memoized("team")
				assert.NoError(t, err)
				assert.Equal(t, []client.Team{{Name: "team"}}, teams)
			})
		}
		wg.Wait()

		assert.Equal(t, int32(1), callCount.Load())
	})
}
