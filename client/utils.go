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

// Helper function to convert any to map[string]any
func convertToMap(data any) (map[string]any, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Helper function to convert map[string]any back to struct
func convertMapBackToStruct(data map[string]any, target any) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonBytes, target)
}
