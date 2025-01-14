package env0

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/stretchr/testify/assert"
)

func TestValidateUrl(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		expectError bool
	}{
		{
			name:        "valid http url",
			url:         "http://example.com",
			expectError: false,
		},
		{
			name:        "valid https url",
			url:         "https://example.com",
			expectError: false,
		},
		{
			name:        "valid url with path",
			url:         "https://example.com/path/to/something",
			expectError: false,
		},
		{
			name:        "invalid url - missing protocol",
			url:         "example.com",
			expectError: true,
		},
		{
			name:        "invalid url - malformed",
			url:         "not-a-url",
			expectError: true,
		},
		{
			name:        "empty url",
			url:         "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := ValidateUrl(tt.url, cty.Path{})
			hasError := diags.HasError()
			assert.Equal(t, tt.expectError, hasError)
		})
	}
}

func TestNewStringInValidator(t *testing.T) {
	allowedValues := []string{"one", "two", "three"}
	validator := NewStringInValidator(allowedValues)

	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{
			name:        "valid value",
			value:       "one",
			expectError: false,
		},
		{
			name:        "another valid value",
			value:       "two",
			expectError: false,
		},
		{
			name:        "invalid value",
			value:       "four",
			expectError: true,
		},
		{
			name:        "empty string",
			value:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := validator(tt.value, cty.Path{})
			hasError := diags.HasError()
			assert.Equal(t, tt.expectError, hasError)
		})
	}
}

func TestValidateNotEmptyString(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{
			name:        "non-empty string",
			value:       "valid",
			expectError: false,
		},
		{
			name:        "empty string",
			value:       "",
			expectError: true,
		},
		{
			name:        "whitespace string",
			value:       "   ",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := ValidateNotEmptyString(tt.value, cty.Path{})
			hasError := diags.HasError()
			assert.Equal(t, tt.expectError, hasError)
		})
	}
}

func TestValidateRetries(t *testing.T) {
	tests := []struct {
		name        string
		value       int
		expectError bool
	}{
		{
			name:        "valid retries - 1",
			value:       1,
			expectError: false,
		},
		{
			name:        "valid retries - 2",
			value:       2,
			expectError: false,
		},
		{
			name:        "valid retries - 3",
			value:       3,
			expectError: false,
		},
		{
			name:        "invalid retries - too low",
			value:       0,
			expectError: true,
		},
		{
			name:        "invalid retries - too high",
			value:       4,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := ValidateRetries(tt.value, cty.Path{})
			hasError := diags.HasError()
			assert.Equal(t, tt.expectError, hasError)
		})
	}
}
