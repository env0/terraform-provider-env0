package client

//go:generate mockgen -destination=project_mock.go -package=client . ProjectApiClient

type ProjectApiClient interface {
	Projects() ([]Project, error)
	Project(id string) (Project, error)
	ProjectCreate(name string) (Project, error)
	ProjectDelete(id string) error
}

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

func (self *ApiClient) ProjectCreate(name string, description string) (Project, error) {
	var result Project
	organizationId, err := self.organizationId()
	if err != nil {
		return Project{}, err
	}

	request := map[string]interface{}{"name": name, "organizationId": organizationId, "description": description}
	err = self.http.Post("/projects", request, &result)
	if err != nil {
		return Project{}, err
	}
	return result, nil
}

func (self *ApiClient) ProjectDelete(id string) error {
	return self.http.Delete("/projects/" + id)
}

func (self *ApiClient) ProjectUpdate(id string, payload UpdateProjectPayload) (Project, error) {
	var result Project
	err := self.http.Put("/projects/"+id, payload, &result)

	if err != nil {
		return Project{}, err
	}
	return result, nil
}
