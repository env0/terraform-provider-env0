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

func TestValidateTags(t *testing.T) {
	tests := []struct {
		name        string
		tags        map[string]any
		expectError bool
	}{
		{
			name:        "single value",
			tags:        map[string]any{"owner": "alice"},
			expectError: false,
		},
		{
			name:        "multi value joined with no space",
			tags:        map[string]any{"team": "eng,payments"},
			expectError: false,
		},
		{
			name:        "space inside a value is legal",
			tags:        map[string]any{"team": "platform eng,payments"},
			expectError: false,
		},
		{
			name:        "space after the separator",
			tags:        map[string]any{"team": "eng, payments"},
			expectError: true,
		},
		{
			name:        "space before the separator",
			tags:        map[string]any{"team": "eng ,payments"},
			expectError: true,
		},
		{
			name:        "empty value",
			tags:        map[string]any{"team": ""},
			expectError: true,
		},
		{
			// legal server-side, so it can arrive by import - but it silently stops matching a filter on "alice"
			name:        "padded single value",
			tags:        map[string]any{"owner": " alice"},
			expectError: true,
		},
		{
			name:        "trailing space on a single value",
			tags:        map[string]any{"owner": "alice "},
			expectError: true,
		},
		{
			name:        "empty segment",
			tags:        map[string]any{"team": "eng,,payments"},
			expectError: true,
		},
		{
			name:        "space inside a key is legal",
			tags:        map[string]any{"cost center": "platform"},
			expectError: false,
		},
		{
			// legal server-side, so it can arrive by import - but the key silently stops matching a filter on "owner"
			name:        "padded key",
			tags:        map[string]any{" owner": "alice"},
			expectError: true,
		},
		{
			name:        "trailing space on a key",
			tags:        map[string]any{"owner ": "alice"},
			expectError: true,
		},
		{
			name:        "no tags",
			tags:        map[string]any{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := ValidateTags(tt.tags, cty.Path{})
			assert.Equal(t, tt.expectError, diags.HasError())
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

func TestNewTerraformVersionValidator(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{name: "version", value: "1.5.7", expectError: false},
		{name: "lowest allowed by the code", value: "RESOLVE_FROM_TERRAFORM_CODE", expectError: false},
		{name: "newest allowed by the code", value: "LATEST_ALLOWED_FROM_CODE", expectError: false},
		{name: "latest", value: "latest", expectError: false},
		{name: "the OpenTofu keyword", value: "RESOLVE_FROM_CODE", expectError: true},
		{name: "a partial keyword", value: "LATEST_ALLOWED", expectError: true},
		{name: "not a version", value: "1.x", expectError: true},
		{name: "a version with a suffix", value: "1.5.7abc", expectError: true},
		{name: "a pre-release version", value: "1.5.0-rc1", expectError: true},
		{name: "a keyword with a prefix", value: "xlatest", expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := NewTerraformVersionValidator()(tt.value, cty.Path{})
			assert.Equal(t, tt.expectError, diags.HasError())
		})
	}
}

func TestNewOpenTofuVersionValidator(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{name: "version", value: "1.13.9", expectError: false},
		{name: "pre-release version", value: "1.13.0-rc1", expectError: false},
		{name: "lowest allowed by the code", value: "RESOLVE_FROM_CODE", expectError: false},
		{name: "newest allowed by the code", value: "LATEST_ALLOWED_FROM_CODE", expectError: false},
		{name: "latest", value: "latest", expectError: false},
		{name: "the Terraform keyword", value: "RESOLVE_FROM_TERRAFORM_CODE", expectError: true},
		{name: "a partial keyword", value: "LATEST_ALLOWED", expectError: true},
		{name: "not a version", value: "1.x", expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := NewOpenTofuVersionValidator()(tt.value, cty.Path{})
			assert.Equal(t, tt.expectError, diags.HasError())
		})
	}
}
