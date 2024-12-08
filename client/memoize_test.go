package client_test

import (
	"errors"
	"testing"

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
}
