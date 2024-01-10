package client

type ModuleSshKey struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Module struct {
	ModuleName           string         `json:"moduleName"`
	ModuleProvider       string         `json:"moduleProvider"`
	Repository           string         `json:"repository"`
	Description          string         `json:"description"`
	LogoUrl              string         `json:"logoUrl"`
	TokenId              string         `json:"tokenId"`
	TokenName            string         `json:"tokenName"`
	GithubInstallationId *int           `json:"githubInstallationId" tfschema:",omitempty"`
	BitbucketClientKey   *string        `json:"bitbucketClientKey" tfschema:",omitempty"`
	IsGitlab             bool           `json:"isGitLab"`
	SshKeys              []ModuleSshKey `json:"sshkeys"`
	Type                 string         `json:"type"`
	Id                   string         `json:"id"`
	OrganizationId       string         `json:"organizationId"`
	Author               User           `json:"author"`
	AuthorId             string         `json:"authorId"`
	CreatedAt            string         `json:"createdAt"`
	UpdatedAt            string         `json:"updatedAt"`
	IsDeleted            bool           `json:"isDeleted"`
	Path                 string         `json:"path"`
	TagPrefix            string         `json:"tagPrefix"`
}

type ModuleCreatePayload struct {
	ModuleName           string         `json:"moduleName"`
	ModuleProvider       string         `json:"moduleProvider"`
	Repository           string         `json:"repository"`
	Description          string         `json:"description,omitempty"`
	LogoUrl              string         `json:"logoUrl,omitempty"`
	TokenId              string         `json:"tokenId,omitempty"`
	TokenName            string         `json:"tokenName,omitempty"`
	GithubInstallationId *int           `json:"githubInstallationId,omitempty"`
	BitbucketClientKey   string         `json:"bitbucketClientKey,omitempty"`
	IsGitlab             *bool          `json:"isGitLab,omitempty"`
	SshKeys              []ModuleSshKey `json:"sshkeys,omitempty"`
	Path                 string         `json:"path,omitempty"`
	TagPrefix            string         `json:"tagPrefix,omitempty"`
}

type ModuleCreatePayloadWith struct {
	ModuleCreatePayload
	Type           string `json:"type"`
	OrganizationId string `json:"organizationId"`
}

type ModuleUpdatePayload struct {
	ModuleName           string         `json:"moduleName,omitempty"`
	ModuleProvider       string         `json:"moduleProvider,omitempty"`
	Repository           string         `json:"repository,omitempty"`
	Description          string         `json:"description,omitempty"`
	LogoUrl              string         `json:"logoUrl,omitempty"`
	TokenId              string         `json:"tokenId"`
	TokenName            string         `json:"tokenName"`
	GithubInstallationId *int           `json:"githubInstallationId"`
	BitbucketClientKey   string         `json:"bitbucketClientKey"`
	IsGitlab             bool           `json:"isGitLab"`
	SshKeys              []ModuleSshKey `json:"sshkeys"`
	Path                 string         `json:"path"`
	TagPrefix            string         `json:"tagPrefix,omitempty"`
}

func (client *ApiClient) ModuleCreate(payload ModuleCreatePayload) (*Module, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	payloadWith := ModuleCreatePayloadWith{
		ModuleCreatePayload: payload,
		OrganizationId:      organizationId,
		Type:                "module",
	}

	var result Module
	if err := client.http.Post("/modules", payloadWith, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) Module(id string) (*Module, error) {
	var result Module
	if err := client.http.Get("/modules/"+id, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) ModuleDelete(id string) error {
	return client.http.Delete("/modules/"+id, nil)
}

func (client *ApiClient) ModuleUpdate(id string, payload ModuleUpdatePayload) (*Module, error) {
	var result Module
	if err := client.http.Patch("/modules/"+id, payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) Modules() ([]Module, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	var result []Module
	if err := client.http.Get("/modules", map[string]string{"organizationId": organizationId}, &result); err != nil {
		return nil, err
	}

	return result, err
}
