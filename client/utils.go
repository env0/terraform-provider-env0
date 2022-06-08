package client

import "encoding/json"

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
