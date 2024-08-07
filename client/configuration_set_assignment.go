package client

import (
	"fmt"
	"strings"
)

func (client *ApiClient) AssignConfigurationSets(scope string, scopeId string, sets []string) error {
	setIds := strings.Join(sets, ",")
	url := fmt.Sprintf("/configuration-sets/assignments/%s/%s?setIds=%s", scope, scopeId, setIds)

	return client.http.Post(url, nil, nil)
}

func (client *ApiClient) UnassignConfigurationSets(scope string, scopeId string, sets []string) error {
	setIds := strings.Join(sets, ",")
	url := fmt.Sprintf("/configuration-sets/assignments/%s/%s", scope, scopeId)

	return client.http.Delete(url, map[string]string{"setIds": setIds})
}

func (client *ApiClient) ConfigurationSetsAssignments(scope string, scopeId string) ([]ConfigurationSet, error) {
	var result []ConfigurationSet

	url := fmt.Sprintf("/configuration-sets/assignments/%s/%s", strings.ToLower(scope), scopeId)

	if err := client.http.Get(url, nil, &result); err != nil {
		return nil, err
	}

	return result, nil
}
