package client

func (self *ApiClient) Projects() ([]Project, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return nil, err
	}
	var result []Project
	err = self.http.Get("/projects", map[string]string{"organizationId": organizationId}, &result)
	if err != nil {
		return []Project{}, err
	}
	return result, nil
}

func (self *ApiClient) Project(id string) (Project, error) {
	var result Project
	err := self.http.Get("/projects/"+id, nil, &result)
	if err != nil {
		return Project{}, err
	}
	return result, nil
}

func (self *ApiClient) ProjectCreate(payload ProjectCreatePayload) (Project, error) {
	var result Project
	organizationId, err := self.organizationId()
	if err != nil {
		return Project{}, err
	}

	request := map[string]interface{}{"name": payload.Name, "organizationId": organizationId, "description": payload.Description}
	err = self.http.Post("/projects", request, &result)
	if err != nil {
		return Project{}, err
	}
	return result, nil
}

func (self *ApiClient) ProjectDelete(id string) error {
	return self.http.Delete("/projects/" + id)
}

func (self *ApiClient) ProjectUpdate(id string, payload ProjectCreatePayload) (Project, error) {
	var result Project
	err := self.http.Put("/projects/"+id, payload, &result)

	if err != nil {
		return Project{}, err
	}
	return result, nil
}
