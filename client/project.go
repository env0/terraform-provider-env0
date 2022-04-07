package client

type Project struct {
	IsArchived     bool   `json:"isArchived"`
	OrganizationId string `json:"organizationId"`
	UpdatedAt      string `json:"updatedAt"`
	CreatedAt      string `json:"createdAt"`
	Id             string `json:"id"`
	Name           string `json:"name"`
	CreatedBy      string `json:"createdBy"`
	Role           string `json:"role"`
	CreatedByUser  User   `json:"createdByUser"`
	Description    string `json:"description"`
}

type ProjectCreatePayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (client *ApiClient) Projects() ([]Project, error) {
	organizationId, err := client.organizationId()
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
	organizationId, err := client.organizationId()
	if err != nil {
		return Project{}, err
	}

	request := map[string]interface{}{"name": payload.Name, "organizationId": organizationId, "description": payload.Description}
	err = client.http.Post("/projects", request, &result)
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
