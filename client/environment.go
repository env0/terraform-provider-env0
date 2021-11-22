package client

func (self *ApiClient) Environments() ([]Environment, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return nil, err
	}
	var result []Environment
	err = self.http.Get("/environments", map[string]string{"organizationId": organizationId}, &result)
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

func (self *ApiClient) EnvironmentCreate(payload EnvironmentCreatePayload) (Environment, error) {
	var result Environment
	request := map[string]interface{}{"name": payload.Name, "projectId": payload.ProjectId}
	err := self.http.Post("/environments", request, &result)
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

func (self *ApiClient) EnvironmentUpdate(id string, payload EnvironmentUpdatePayload) (Environment, error) {
	var result Environment
	err := self.http.Put("/projects/"+id, payload, &result)

	if err != nil {
		return Environment{}, err
	}
	return result, nil
}
