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

	result, err := self.http.Post("/ssh-keys", payload)
	if err != nil {
		return SshKey{}, err
	}
	return result.(SshKey), nil
}

func (self *ApiClient) SshKeyDelete(id string) error {
	return self.http.Delete("/ssh-keys/" + id)
}

func (self *ApiClient) SshKeys() ([]SshKey, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return nil, err
	}

	result, err := self.http.Get("/ssh-keys", map[string]string{"organizationId": organizationId})
	if err != nil {
		return nil, err
	}

	sshKeys := result.([]SshKey)
	return sshKeys, err
}
