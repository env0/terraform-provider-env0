package client

func (self *ApiClient) AwsCredentialsCreate(request AwsCredentialsCreatePayload) (ApiKey, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return ApiKey{}, err
	}

	request.Type = "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT"
	request.OrganizationId = organizationId

	result, err := self.http.Post("/credentials", request)
	if err != nil {
		return ApiKey{}, err
	}
	return result.(ApiKey), nil
}