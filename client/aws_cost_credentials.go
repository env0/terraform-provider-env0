package client

import (
	"fmt"
)

func (self *ApiClient) AwsCostCredentials(id string) (ApiKey, error) {
	var credentials, err = self.CloudCredentialsList()
	if err != nil {
		return ApiKey{}, err
	}

	for _, v := range credentials {
		if v.Id == id {
			return v, nil
		}
	}

	return ApiKey{}, fmt.Errorf("AwsCostCredentials: [%s] not found ", id)
}

func (self *ApiClient) AwsCostCredentialsList() ([]ApiKey, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return []ApiKey{}, err
	}

	var credentials []ApiKey
	err = self.http.Get("/credentials", map[string]string{"organizationId": organizationId}, &credentials)
	if err != nil {
		return []ApiKey{}, err
	}

	return credentials, nil
}

func (self *ApiClient) AwsCostCredentialsCreate(request AwsCredentialsCreatePayload) (ApiKey, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return ApiKey{}, err
	}

	request.Type = "AWS_ASSUMED_ROLE"
	request.OrganizationId = organizationId

	var result ApiKey
	err = self.http.Post("/credentials", request, &result)
	if err != nil {
		return ApiKey{}, err
	}
	return result, nil
}

func (self *ApiClient) AwsCostCredentialsDelete(id string) error {
	return self.http.Delete("/credentials/" + id)
}
