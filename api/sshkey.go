package api

import (
	"errors"
	"net/url"
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
	err = self.postJSON("/ssh-keys", payload, &result)
	if err != nil {
		return SshKey{}, err
	}
	return result, nil
}

func (self *ApiClient) SshKeyDelete(id string) error {
	return self.delete("/ssh-keys/" + id)
}

func (self *ApiClient) SshKeys() ([]SshKey, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return nil, err
	}
	var result []SshKey
	params := url.Values{}
	params.Add("organizationId", organizationId)
	err = self.getJSON("/ssh-keys", params, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}
