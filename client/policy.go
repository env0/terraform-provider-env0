package client

import "errors"

func (self *ApiClient) Policy() (Policy, error) {
	var result []Policy
	err := self.http.Get("/policies", nil, &result)
	if err != nil {
		return Policy{}, err
	}
	if len(result) != 1 {
		return Policy{}, errors.New("Server responded with too many policies")
	}
	return result[0], nil
}

func (self *ApiClient) PolicyUpdate(id string, payload PolicyUpdatePayload) (Policy, error) {
	if payload.ProjectId == "" {
		return Policy{}, errors.New("Must specify project ID on update")
	}

	var result Policy
	if err := self.http.Put("/policies/"+id, payload, &result); err != nil {
		return Policy{}, err
	}
	return result, nil
}
