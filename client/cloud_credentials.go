package client

import (
	"fmt"
)

func (self *ApiClient) CloudCredentials(id string) (Credentials, error) {
	var credentials, err = self.CloudCredentialsList()
	if err != nil {
		return Credentials{}, err
	}

	for _, v := range credentials {
		if v.Id == id {
			return v, nil
		}
	}

	return Credentials{}, fmt.Errorf("CloudCredentials: [%s] not found ", id)
}

func (self *ApiClient) CloudCredentialsList() ([]Credentials, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return []Credentials{}, err
	}

	var credentials []Credentials
	err = self.http.Get("/credentials", map[string]string{"organizationId": organizationId}, &credentials)
	if err != nil {
		return []Credentials{}, err
	}

	return credentials, nil
}

func (self *ApiClient) AwsCredentialsCreate(request AwsCredentialsCreatePayload) (Credentials, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return Credentials{}, err
	}

	request.OrganizationId = organizationId

	var result Credentials
	err = self.http.Post("/credentials", request, &result)
	if err != nil {
		return Credentials{}, err
	}
	return result, nil
}

func (self *ApiClient) GcpCredentialsCreate(request GcpCredentialsCreatePayload) (Credentials, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return Credentials{}, err
	}

	request.OrganizationId = organizationId

	var result Credentials
	err = self.http.Post("/credentials", request, &result)
	if err != nil {
		return Credentials{}, err
	}
	return result, nil
}

func (self *ApiClient) AzureCredentialsCreate(request AzureCredentialsCreatePayload) (Credentials, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return Credentials{}, err
	}

	request.OrganizationId = organizationId

	var result Credentials
	err = self.http.Post("/credentials", request, &result)
	if err != nil {
		return Credentials{}, err
	}
	return result, nil
}

func (self *ApiClient) CloudCredentialsDelete(id string) error {
	return self.http.Delete("/credentials/" + id)
}

func (self *ApiClient) GoogleCostCredentialsCreate(request GoogleCostCredentialsCreatePayload) (Credentials, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return Credentials{}, err
	}

	request.OrganizationId = organizationId
	var result Credentials
	err = self.http.Post("/credentials", request, &result)
	if err != nil {
		return Credentials{}, err
	}
	return result, nil
}
