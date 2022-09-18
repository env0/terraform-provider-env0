package client

type RoleCreatePayload struct {
	Name           string   `json:"name"`
	OrganizationId string   `json:"organizationId"`
	Permissions    []string `json:"permissions"`
	IsDefaultRole  bool     `json:"isDefaultRole"`
}

type RoleUpdatePayload struct {
	Name          string   `json:"name"`
	Permissions   []string `json:"permissions"`
	IsDefaultRole bool     `json:"isDefaultRole"`
}

type Role struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	OrganizationId string   `json:"organizationId"`
	Permissions    []string `json:"permissions"`
	IsDefaultRole  bool     `json:"isDefaultRole"`
}

func (client *ApiClient) RoleCreate(payload RoleCreatePayload) (*Role, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	payload.OrganizationId = organizationId

	var result Role
	if err := client.http.Post("/roles", payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) Role(id string) (*Role, error) {
	var result Role

	if err := client.http.Get("/roles/"+id, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) RoleDelete(id string) error {
	return client.http.Delete("/roles/" + id)
}

func (client *ApiClient) RoleUpdate(id string, payload RoleUpdatePayload) (*Role, error) {
	var result Role

	if err := client.http.Put("/roles/"+id, payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) Roles() ([]Role, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	var result []Role
	if err := client.http.Get("/roles", map[string]string{"organizationId": organizationId}, &result); err != nil {
		return nil, err
	}

	return result, err
}
