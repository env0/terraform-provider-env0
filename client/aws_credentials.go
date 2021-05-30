package client

import (
	"fmt"
)

func (self *ApiClient) AwsCredentials(id string) (ApiKey, error) {
	var credentials, err = self.AwsCredentialsList()
	if err != nil {
		return ApiKey{}, err
	}

	for _, v := range credentials {
		if v.Id == id {
			return v, nil
		}
	}

	return ApiKey{}, fmt.Errorf("AwsCredentials: [%s] not found ", id)
}

func (self *ApiClient) AwsCredentialsList() ([]ApiKey, error) {
	var credentials []ApiKey
	err := self.http.Get("/credentials", nil, &credentials)
	if err != nil {
		return []ApiKey{}, err
	}

	return credentials, nil
}

func (self *ApiClient) AwsCredentialsCreate(request AwsCredentialsCreatePayload) (ApiKey, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return ApiKey{}, err
	}

	request.Type = "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT"
	request.OrganizationId = organizationId

	var result ApiKey
	err = self.http.Post("/credentials", request, &result)
	if err != nil {
		return ApiKey{}, err
	}
	return result, nil
}

func (self *ApiClient) AwsCredentialsDelete(id string) error {
	return self.http.Delete("/credentials/" + id)
}