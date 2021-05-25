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

	result, err := self.http.Get("/projects", map[string]string{"organizationId": organizationId})
	if err != nil {
		return []Project{}, err
	}
	return result.([]Project), nil
}

func (self *ApiClient) Project(id string) (Project, error) {
	result, err := self.http.Get("/projects/"+id, nil)
	if err != nil {
		return Project{}, err
	}
	return result.(Project), nil
}

func (self *ApiClient) ProjectCreate(name string) (Project, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return Project{}, err
	}

	request := map[string]interface{}{"name": name, "organizationId": organizationId}
	result, err := self.http.Post("/projects", request)
	if err != nil {
		return Project{}, err
	}
	return result.(Project), nil
}

func (self *ApiClient) ProjectDelete(id string) error {
	return self.http.Delete("/projects/" + id)
}
