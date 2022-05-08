package client

type OrganizationUser struct {
	User   User   `json:"user"`
	Role   string `json:"role"`
	Status string `json:"status"`
}

func (client *ApiClient) Users() ([]OrganizationUser, error) {
	organizationId, err := client.organizationId()
	if err != nil {
		return nil, err
	}

	var result []OrganizationUser
	if err := client.http.Get("/organizations/"+organizationId+"/users", nil, &result); err != nil {
		return nil, err
	}

	return result, nil
}
