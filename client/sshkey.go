package client

import (
	"errors"
)

func (self *ApiClient) SshKeyCreate(payload SshKeyCreatePayload) (SshKey, error) {
	if payload.Name == "" {
		return SshKey{}, errors.New("Must specify ssh key name on creation")
	}
	if payload.Value == "" {
		return SshKey{}, errors.New("Must specify ssh key value (private key in PEM format) on creation")
	}
	if payload.OrganizationId != "" {
		return SshKey{}, errors.New("Must not specify organizationId")
	}
	organizationId, err := self.organizationId()
	if err != nil {
		return SshKey{}, nil
	}
	payload.OrganizationId = organizationId

	var result SshKey
	err = self.http.Post("/ssh-keys", payload, &result)
	if err != nil {
		return SshKey{}, err
	}
	return result, nil
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
