package client

import "fmt"

// Policy retrieves a policy from the API
func (self *ApiClient) Policy(projectId string) (Policy, error) {
	u, err := newQueryURL("/policies", parameter{"projectId", projectId})
	if err != nil {
		return Policy{}, err
	}

	var result Policy
	err = self.http.Get(u.String(), nil, &result)
	if err != nil {
		return Policy{}, err
	}
	return result, nil
}

// PolicyUpdate updates a policy through the API
func (self *ApiClient) PolicyUpdate(payload PolicyUpdatePayload) (Policy, error) {
	var result Policy
	if payload.ProjectId == "" {
		return Policy{}, fmt.Errorf("projectId is required when updating policy")
	}
	err := self.http.Put("/policies", payload, &result)
	if err != nil {
		return Policy{}, err
	}
	return result, nil
}
