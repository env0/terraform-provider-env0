package client

// Policy retrieves a policy from the API
func (self *ApiClient) Policy() (Policy, error) {
	var result Policy
	err := self.http.Get("/policies", nil, &result)
	if err != nil {
		return Policy{}, err
	}
	return result, nil
}

// PolicyUpdate updates a policy through the API
func (self *ApiClient) PolicyUpdate(id string, payload PolicyUpdatePayload) (Policy, error) {
	var result Policy
	if err := self.http.Put("/policies/"+id, payload, &result); err != nil {
		return Policy{}, err
	}
	return result, nil
}
