package client

import "net/url"

// Policy retrieves a policy from the API
func (self *ApiClient) Policy(projectId string) (Policy, error) {
	u, err := url.Parse("/policies")
	if err != nil {
		return Policy{}, err
	}

	q := u.Query()
	q.Add("projectId", projectId)
	u.RawQuery = q.Encode()

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
	err := self.http.Put("/policies", payload, &result)
	if err != nil {
		return Policy{}, err
	}
	return result, nil
}
