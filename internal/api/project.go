package api

func (self *ApiClient) Projects() ([]Project, error) {
	organizationId, err := self.getOrganizationId()
	if err != nil {
		return nil, err
	}
	var result []Project
	err = self.client.Get("/projects", map[string]string{"organizationId": organizationId}, &result)
	if err != nil {
		return []Project{}, err
	}
	return result, nil
}

func (self *ApiClient) Project(id string) (Project, error) {
	var result Project
	err := self.client.Get("/projects/"+id, nil, &result)
	if err != nil {
		return Project{}, err
	}
	return result, nil
}

func (self *ApiClient) ProjectCreate(name string) (Project, error) {
	var result Project
	request := map[string]interface{}{"name": name}
	err := self.client.Post("/projects", request, &result)
	if err != nil {
		return Project{}, err
	}
	return result, nil
}

func (self *ApiClient) ProjectDelete(id string) error {
	return self.client.Delete("/projects/" + id)
}
