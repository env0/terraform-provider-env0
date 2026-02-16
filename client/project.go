package client

type Project struct {
	IsArchived      bool   `json:"isArchived"`
	OrganizationId  string `json:"organizationId"`
	UpdatedAt       string `json:"updatedAt"`
	CreatedAt       string `json:"createdAt"`
	Id              string `json:"id"`
	Name            string `json:"name"`
	CreatedBy       string `json:"createdBy"`
	Role            string `json:"role"`
	CreatedByUser   User   `json:"createdByUser"`
	Description     string `json:"description"`
	ParentProjectId string `json:"parentProjectId,omitempty" tfschema:",omitempty"`
	Hierarchy       string `json:"hierarchy"`
}

type ProjectCreatePayload struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	ParentProjectId string `json:"parentProjectId,omitempty"`
}

type ProjectUpdatePayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ModuleTestingProject struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func (client *ApiClient) Projects() ([]Project, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	var result []Project

	if err := client.http.Get("/projects", map[string]string{"organizationId": organizationId}, &result); err != nil {
		return []Project{}, err
	}

	return result, nil
}

func (client *ApiClient) Project(id string) (Project, error) {
	var result Project

	if err := client.http.Get("/projects/"+id, nil, &result); err != nil {
		return Project{}, err
	}

	return result, nil
}

func (client *ApiClient) ProjectCreate(payload ProjectCreatePayload) (Project, error) {
	var result Project

	organizationId, err := client.OrganizationId()
	if err != nil {
		return Project{}, err
	}

	payloadWithOrganizationId := struct {
		ProjectCreatePayload
		OrganizationId string `json:"organizationId"`
	}{
		payload,
		organizationId,
	}

	err = client.http.Post("/projects", payloadWithOrganizationId, &result)
	if err != nil {
		return Project{}, err
	}

	return result, nil
}

func (client *ApiClient) ProjectDelete(id string) error {
	return client.http.Delete("/projects/"+id, nil)
}

func (client *ApiClient) ProjectUpdate(id string, payload ProjectUpdatePayload) (Project, error) {
	var result Project

	err := client.http.Put("/projects/"+id, payload, &result)
	if err != nil {
		return Project{}, err
	}

	return result, nil
}

func (client *ApiClient) ProjectMove(id string, targetProjectId string) error {
	// Pass nil if a subproject becomes a project.
	var targetProjectIdPtr *string
	if targetProjectId != "" {
		targetProjectIdPtr = &targetProjectId
	}

	payload := struct {
		TargetProjectId *string `json:"targetProjectId"`
	}{
		targetProjectIdPtr,
	}

	return client.http.Post("/projects/"+id+"/move", payload, nil)
}

func (client *ApiClient) ModuleTestingProject() (*ModuleTestingProject, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	var result ModuleTestingProject
	if err := client.http.Get("/projects/modules/testing/"+organizationId, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
