package client

func (self *ApiClient) SshKeyCreate(payload SshKeyCreatePayload) (*SshKey, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return nil, err
	}
	payload.OrganizationId = organizationId

	var result SshKey
	if err := self.http.Post("/ssh-keys", payload, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (self *ApiClient) SshKeyDelete(id string) error {
	return self.http.Delete("/ssh-keys/" + id)
}

func (self *ApiClient) SshKeys() ([]SshKey, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return nil, err
	}
	var result []SshKey
	err = self.http.Get("/ssh-keys", map[string]string{"organizationId": organizationId}, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}
