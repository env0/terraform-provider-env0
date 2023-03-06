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
}

type ProjectCreatePayload struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	ParentProjectId string `json:"parentProjectId,omitempty"`
}

func (client *ApiClient) Projects() ([]Project, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}
	var result []Project
	err = client.http.Get("/projects", map[string]string{"organizationId": organizationId}, &result)
	if err != nil {
		return []Project{}, err
	}
	return result, nil
}

func (client *ApiClient) Project(id string) (Project, error) {
	var result Project
	err := client.http.Get("/projects/"+id, nil, &result)
	if err != nil {
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
	return client.http.Delete("/projects/" + id)
}

func (client *ApiClient) ProjectUpdate(id string, payload ProjectCreatePayload) (Project, error) {
	var result Project
	err := client.http.Put("/projects/"+id, payload, &result)

	if err != nil {
		return Project{}, err
	}
	return result, nil
}
