package client

//templates are actually called "blueprints" in some parts of the API, this layer
//attempts to abstract this detail away - all the users of api client should
//only use "template", no mention of blueprint

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
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

type VariablesFromRepositoryPayload struct {
	BitbucketClientKey   string   `json:"bitbucketClientKey,omitempty"`
	GithubInstallationId int      `json:"githubInstallationId,omitempty"`
	Path                 string   `json:"path"`
	Revision             string   `json:"revision"`
	SshKeyIds            []string `json:"sshKeyIds"`
	TokenId              string   `json:"tokenId,omitempty"`
	Repository           string   `json:"repository"`
}

func (client *ApiClient) TemplateCreate(payload TemplateCreatePayload) (Template, error) {
	if payload.Name == "" {
		return Template{}, errors.New("must specify template name on creation")
	}
	if payload.OrganizationId != "" {
		return Template{}, errors.New("must not specify organizationId")
	}
	if payload.Type != "terragrunt" && payload.TerragruntVersion != "" {
		return Template{}, errors.New("can't define terragrunt version for non-terragrunt blueprint")
	}
	if payload.Type == "terragrunt" && payload.TerragruntVersion == "" {
		return Template{}, errors.New("must supply Terragrunt version")
	}
	organizationId, err := client.organizationId()
	if err != nil {
		return Template{}, nil
	}
	payload.OrganizationId = organizationId

	var result Template
	err = client.http.Post("/blueprints", payload, &result)
	if err != nil {
		return Template{}, err
	}
	return result, nil
}

func (client *ApiClient) Template(id string) (Template, error) {
	var result Template
	err := client.http.Get("/blueprints/"+id, nil, &result)
	if err != nil {
		return Template{}, err
	}
	return result, nil
}

func (client *ApiClient) TemplateDelete(id string) error {
	return client.http.Delete("/blueprints/" + id)
}

func (client *ApiClient) TemplateUpdate(id string, payload TemplateCreatePayload) (Template, error) {
	if payload.Name == "" {
		return Template{}, errors.New("must specify template name on creation")
	}
	if payload.OrganizationId != "" {
		return Template{}, errors.New("must not specify organizationId")
	}
	if payload.Type != "terragrunt" && payload.TerragruntVersion != "" {
		return Template{}, errors.New("can't define terragrunt version for non-terragrunt blueprint")
	}
	if payload.Type == "terragrunt" && payload.TerragruntVersion == "" {
		return Template{}, errors.New("must supply Terragrunt version")
	}
	organizationId, err := client.organizationId()
	if err != nil {
		return Template{}, err
	}
	payload.OrganizationId = organizationId

	var result Template
	err = client.http.Put("/blueprints/"+id, payload, &result)
	if err != nil {
		return Template{}, err
	}
	return result, nil
}

func (client *ApiClient) Templates() ([]Template, error) {
	organizationId, err := client.organizationId()
	if err != nil {
		return nil, err
	}
	var result []Template
	err = client.http.Get("/blueprints", map[string]string{"organizationId": organizationId}, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (client *ApiClient) AssignTemplateToProject(id string, payload TemplateAssignmentToProjectPayload) (Template, error) {
	var result Template
	if payload.ProjectId == "" {
		return result, errors.New("must specify projectId on assignment to a template")
	}
	err := client.http.Patch("/blueprints/"+id+"/projects", payload, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (client *ApiClient) RemoveTemplateFromProject(templateId string, projectId string) error {
	return client.http.Delete("/blueprints/" + templateId + "/projects/" + projectId)
}

func (client *ApiClient) VariablesFromRepository(payload *VariablesFromRepositoryPayload) ([]ConfigurationVariable, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	var paramsInterface map[string]interface{}
	if err := json.Unmarshal(b, &paramsInterface); err != nil {
		return nil, err
	}

	params := map[string]string{}
	for key, value := range paramsInterface {
		if key == "githubInstallationId" {
			params[key] = strconv.Itoa(int(value.(float64)))
		} else if key == "sshKeyIds" {
			sshkeys := []string{}
			if value != nil {
				for _, sshkey := range value.([]interface{}) {
					sshkeys = append(sshkeys, "\""+sshkey.(string)+"\"")
				}
			}
			params[key] = "[" + strings.Join(sshkeys, ",") + "]"
		} else {
			params[key] = value.(string)
		}
	}

	var result []ConfigurationVariable
	if err := client.http.Get("/blueprints/variables-from-repository", params, &result); err != nil {
		return nil, err
	}

	return result, nil
}
