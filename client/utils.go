package client

import "encoding/json"

func toParamsInterface(i any) (map[string]any, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	var paramsInterface map[string]any
	if err := json.Unmarshal(b, &paramsInterface); err != nil {
		return nil, err
	}

	return paramsInterface, nil
}
