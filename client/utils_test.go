package client

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNullString(t *testing.T) {
	st := struct {
		Null    NullString `json:"null"`
		NotNull NullString `json:"notNull"`
	}{
		Null:    "",
		NotNull: "cool",
	}

	b, err := json.Marshal(st)

	require.Nil(t, err)

	s := string(b)

	require.True(t, strings.Contains(s, `"null":null`))
	require.True(t, strings.Contains(s, `"notNull":"cool"`))
}

func TestNullInt(t *testing.T) {
	st := struct {
		Null    NullInt `json:"null"`
		NotNull NullInt `json:"notNull"`
	}{
		Null:    0,
		NotNull: 1000,
	}

	b, err := json.Marshal(st)

	require.Nil(t, err)

	s := string(b)

	require.True(t, strings.Contains(s, `"null":null`))
	require.True(t, strings.Contains(s, `"notNull":1000`))
}
