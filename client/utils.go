package client

import (
	"encoding/json"
	"fmt"
)

// In some cases it's important to pass 'null' explicitly via the API.
// Null* will pass null when it's a "zero" value ("" for string, 0 for int, etc...).

type NullString string

type NullInt int

func (s NullString) MarshalJSON() ([]byte, error) {
	if s == "" {
		return []byte("null"), nil
	}

	return []byte(fmt.Sprintf("\"%s\"", s)), nil
}

func (i NullInt) MarshalJSON() ([]byte, error) {
	if i == 0 {
		return []byte("null"), nil
	}

	return []byte(fmt.Sprintf("%d", i)), nil
}

func toParamsInterface(i interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	var paramsInterface map[string]interface{}
	if err := json.Unmarshal(b, &paramsInterface); err != nil {
		return nil, err
	}

	return paramsInterface, nil
}
