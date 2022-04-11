package client

type SshKey struct {
	User           User   `json:"user"`
	UserId         string `json:"userId"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
	Id             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	OrganizationId string `json:"organizationId"`
	Value          string `json:"value"`
}

type SshKeyCreatePayload struct {
	Name           string `json:"name"`
	OrganizationId string `json:"organizationId"`
	Value          string `json:"value"`
}

func (client *ApiClient) SshKeyCreate(payload SshKeyCreatePayload) (*SshKey, error) {
	organizationId, err := client.organizationId()
	if err != nil {
		return nil, err
	}
	payload.OrganizationId = organizationId

	var result SshKey
	if err := client.http.Post("/ssh-keys", payload, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *ApiClient) SshKeyDelete(id string) error {
	return client.http.Delete("/ssh-keys/" + id)
}

func (client *ApiClient) SshKeys() ([]SshKey, error) {
	organizationId, err := client.organizationId()
	if err != nil {
		return nil, err
	}
	var result []SshKey
	err = client.http.Get("/ssh-keys", map[string]string{"organizationId": organizationId}, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}
