package client

type VscToken struct {
	Token int `json:"token"`
}

func (client *ApiClient) VcsToken(vcsType string, repository string) (*VscToken, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	var result VscToken
	if err := client.http.Get("/vcs-token/"+vcsType, map[string]string{
		"organizationId": organizationId,
		"repository":     repository,
	}, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
