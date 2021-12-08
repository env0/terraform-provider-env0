package client

func (self *ApiClient) Environments() ([]Environment, error) {
	var result []Environment
	err := self.http.Get("/environments", nil, &result)
	if err != nil {
		return []Environment{}, err
	}
	return result, nil
}

func (self *ApiClient) Environment(id string) (Environment, error) {
	var result Environment
	err := self.http.Get("/environments/"+id, nil, &result)
	if err != nil {
		return Environment{}, err
	}
	return result, nil
}

func (self *ApiClient) EnvironmentCreate(payload EnvironmentCreate) (Environment, error) {
	var result Environment

	err := self.http.Post("/environments", payload, &result)
	if err != nil {
		return Environment{}, err
	}
	return result, nil
}

func (self *ApiClient) EnvironmentDestroy(id string) (Environment, error) {
	var result Environment
	err := self.http.Post("/environments/"+id+"/destroy", nil, &result)
	if err != nil {
		return Environment{}, err
	}
	return result, nil
}

func (self *ApiClient) EnvironmentUpdate(id string, payload EnvironmentUpdate) (Environment, error) {
	var result Environment
	err := self.http.Put("/environments/"+id, payload, &result)

	if err != nil {
		return Environment{}, err
	}
	return result, nil
}

func (self *ApiClient) EnvironmentUpdateTTL(id string, payload TTL) (Environment, error) {
	var result Environment
	err := self.http.Put("/environments/"+id+"/ttl", payload, &result)

	if err != nil {
		return Environment{}, err
	}
	return result, nil
}

func (self *ApiClient) EnvironmentDeploy(id string, payload DeployRequest) (EnvironmentDeployResponse, error) {
	var result EnvironmentDeployResponse
	err := self.http.Post("/environments/"+id+"/deployments", payload, &result)

	if err != nil {
		return EnvironmentDeployResponse{}, err
	}
	return result, nil
}
