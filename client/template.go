package client

//templates are actually called "blueprints" in some parts of the API, this layer
//attempts to abstract this detail away - all the users of api client should
//only use "template", no mention of blueprint

import (
	"errors"
)

type TemplateRetryOn struct {
	Times      int    `json:"times,omitempty"`
	ErrorRegex string `json:"errorRegex"`
}

type TemplateRetry struct {
	OnDeploy  *TemplateRetryOn `json:"onDeploy"`
	OnDestroy *TemplateRetryOn `json:"onDestroy"`
}

type TemplateType string

const (
	TemplateTypeTerraform  TemplateType = "terraform"
	TemplateTypeTerragrunt TemplateType = "terragrunt"
)

type TemplateSshKey struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Template struct {
	Author               User             `json:"author"`
	AuthorId             string           `json:"authorId"`
	CreatedAt            string           `json:"createdAt"`
	Href                 string           `json:"href"`
	Id                   string           `json:"id"`
	Name                 string           `json:"name"`
	Description          string           `json:"description"`
	OrganizationId       string           `json:"organizationId"`
	Path                 string           `json:"path"`
	Revision             string           `json:"revision"`
	ProjectId            string           `json:"projectId"`
	ProjectIds           []string         `json:"projectIds"`
	Repository           string           `json:"repository"`
	Retry                TemplateRetry    `json:"retry"`
	SshKeys              []TemplateSshKey `json:"sshKeys"`
	Type                 string           `json:"type"`
	GithubInstallationId int              `json:"githubInstallationId"`
	IsGitlabEnterprise   bool             `json:"isGitLabEnterprise"`
	TokenId              string           `json:"tokenId,omitempty"`
	GitlabProjectId      int              `json:"gitlabProjectId,omitempty"`
	UpdatedAt            string           `json:"updatedAt"`
	TerraformVersion     string           `json:"terraformVersion"`
	TerragruntVersion    string           `json:"terragruntVersion,omitempty"`
	IsDeleted            bool             `json:"isDeleted,omitempty"`
	BitbucketClientKey   string           `json:"bitbucketClientKey"`
	IsGitHubEnterprise   bool             `json:"isGitHubEnterprise"`
	IsBitbucketServer    bool             `json:"isBitbucketServer"`
}

type TemplateCreatePayload struct {
	Retry                TemplateRetry    `json:"retry"`
	SshKeys              []TemplateSshKey `json:"sshKeys,omitempty"`
	Type                 TemplateType     `json:"type"`
	Description          string           `json:"description"`
	Name                 string           `json:"name"`
	Repository           string           `json:"repository"`
	Path                 string           `json:"path"`
	IsGitLab             bool             `json:"isGitLab"`
	TokenName            string           `json:"tokenName"`
	TokenId              string           `json:"tokenId,omitempty"`
	GithubInstallationId int              `json:"githubInstallationId,omitempty"`
	GitlabProjectId      int              `json:"gitlabProjectId,omitempty"`
	Revision             string           `json:"revision"`
	OrganizationId       string           `json:"organizationId"`
	TerraformVersion     string           `json:"terraformVersion"`
	TerragruntVersion    string           `json:"terragruntVersion,omitempty"`
	IsGitlabEnterprise   bool             `json:"isGitLabEnterprise"`
	BitbucketClientKey   string           `json:"bitbucketClientKey,omitempty"`
	IsGitHubEnterprise   bool             `json:"isGitHubEnterprise"`
	IsBitbucketServer    bool             `json:"isBitbucketServer"`
}

type TemplateAssignmentToProjectPayload struct {
	ProjectId string `json:"projectId"`
}

type TemplateAssignmentToProject struct {
	Id         string `json:"id"`
	TemplateId string `json:"templateId"`
	ProjectId  string `json:"projectId"`
}

func (self *ApiClient) TemplateCreate(payload TemplateCreatePayload) (Template, error) {
	if payload.Name == "" {
		return Template{}, errors.New("Must specify template name on creation")
	}
	if payload.OrganizationId != "" {
		return Template{}, errors.New("Must not specify organizationId")
	}
	if payload.Type != "terragrunt" && payload.TerragruntVersion != "" {
		return Template{}, errors.New("Can't define terragrunt version for non-terragrunt blueprint")
	}
	if payload.Type == "terragrunt" && payload.TerragruntVersion == "" {
		return Template{}, errors.New("Must supply Terragrunt version")
	}
	organizationId, err := self.organizationId()
	if err != nil {
		return Template{}, nil
	}
	payload.OrganizationId = organizationId

	var result Template
	err = self.http.Post("/blueprints", payload, &result)
	if err != nil {
		return Template{}, err
	}
	return result, nil
}

func (self *ApiClient) Template(id string) (Template, error) {
	var result Template
	err := self.http.Get("/blueprints/"+id, nil, &result)
	if err != nil {
		return Template{}, err
	}
	return result, nil
}

func (self *ApiClient) TemplateDelete(id string) error {
	return self.http.Delete("/blueprints/" + id)
}

func (self *ApiClient) TemplateUpdate(id string, payload TemplateCreatePayload) (Template, error) {
	if payload.Name == "" {
		return Template{}, errors.New("Must specify template name on creation")
	}
	if payload.OrganizationId != "" {
		return Template{}, errors.New("Must not specify organizationId")
	}
	if payload.Type != "terragrunt" && payload.TerragruntVersion != "" {
		return Template{}, errors.New("Can't define terragrunt version for non-terragrunt blueprint")
	}
	if payload.Type == "terragrunt" && payload.TerragruntVersion == "" {
		return Template{}, errors.New("Must supply Terragrunt version")
	}
	organizationId, err := self.organizationId()
	if err != nil {
		return Template{}, err
	}
	payload.OrganizationId = organizationId

	var result Template
	err = self.http.Put("/blueprints/"+id, payload, &result)
	if err != nil {
		return Template{}, err
	}
	return result, nil
}

func (self *ApiClient) Templates() ([]Template, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return nil, err
	}
	var result []Template
	err = self.http.Get("/blueprints", map[string]string{"organizationId": organizationId}, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (self *ApiClient) AssignTemplateToProject(id string, payload TemplateAssignmentToProjectPayload) (Template, error) {
	var result Template
	if payload.ProjectId == "" {
		return result, errors.New("Must specify projectId on assignment to a template")
	}
	err := self.http.Patch("/blueprints/"+id+"/projects", payload, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (self *ApiClient) RemoveTemplateFromProject(templateId string, projectId string) error {
	return self.http.Delete("/blueprints/" + templateId + "/projects/" + projectId)
}
