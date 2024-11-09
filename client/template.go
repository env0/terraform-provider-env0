package client

// templates are actually called "blueprints" in some parts of the API, this layer
// attempts to abstract this detail away - all the users of api client should
// only use "template", no mention of blueprint

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
)

const TERRAGRUNT = "terragrunt"
const OPENTOFU = "opentofu"
const TERRAFORM = "terraform"
const WORKFLOW = "workflow"

type TemplateRetryOn struct {
	Times      int    `json:"times,omitempty"`
	ErrorRegex string `json:"errorRegex"`
}

type TemplateRetry struct {
	OnDeploy  *TemplateRetryOn `json:"onDeploy"`
	OnDestroy *TemplateRetryOn `json:"onDestroy"`
}

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
	Path                 string           `json:"path" tfschema:",omitempty"`
	Revision             string           `json:"revision"`
	ProjectId            string           `json:"projectId"`
	ProjectIds           []string         `json:"projectIds"`
	Repository           string           `json:"repository"`
	Retry                TemplateRetry    `json:"retry"`
	SshKeys              []TemplateSshKey `json:"sshKeys"`
	Type                 string           `json:"type"`
	GithubInstallationId int              `json:"githubInstallationId" tfschema:",omitempty"`
	IsGitlabEnterprise   bool             `json:"isGitLabEnterprise"`
	TokenId              string           `json:"tokenId" tfschema:",omitempty"`
	UpdatedAt            string           `json:"updatedAt"`
	TerraformVersion     string           `json:"terraformVersion" tfschema:",omitempty"`
	TerragruntVersion    string           `json:"terragruntVersion" tfschema:",omitempty"`
	OpentofuVersion      string           `json:"opentofuVersion" tfschema:",omitempty"`
	IsDeleted            bool             `json:"isDeleted"`
	BitbucketClientKey   string           `json:"bitbucketClientKey" tfschema:",omitempty"`
	IsGithubEnterprise   bool             `json:"isGitHubEnterprise"`
	IsBitbucketServer    bool             `json:"isBitbucketServer"`
	FileName             string           `json:"fileName" tfschema:",omitempty"`
	IsTerragruntRunAll   bool             `json:"isTerragruntRunAll"`
	IsAzureDevOps        bool             `json:"isAzureDevOps" tfschema:"is_azure_devops"`
	IsHelmRepository     bool             `json:"isHelmRepository"`
	HelmChartName        string           `json:"helmChartName" tfschema:",omitempty"`
	IsGitlab             bool             `json:"isGitLab"`
	TerragruntTfBinary   string           `json:"terragruntTfBinary" tfschema:",omitempty"`
	TokenName            string           `json:"tokenName" tfschema:",omitempty"`
	AnsibleVersion       string           `json:"ansibleVersion" tfschema:",omitempty"`
}

type TemplateCreatePayload struct {
	Retry                TemplateRetry    `json:"retry"`
	SshKeys              []TemplateSshKey `json:"sshKeys,omitempty"`
	Type                 string           `json:"type"`
	Description          string           `json:"description"`
	Name                 string           `json:"name"`
	Repository           string           `json:"repository"`
	Path                 string           `json:"path,omitempty"`
	IsGitlab             bool             `json:"isGitLab"`
	TokenName            string           `json:"tokenName,omitempty"`
	TokenId              string           `json:"tokenId,omitempty"`
	GithubInstallationId int              `json:"githubInstallationId,omitempty"`
	Revision             string           `json:"revision"`
	OrganizationId       string           `json:"organizationId"`
	TerraformVersion     string           `json:"terraformVersion,omitempty"`
	TerragruntVersion    string           `json:"terragruntVersion,omitempty"`
	OpentofuVersion      string           `json:"opentofuVersion,omitempty"`
	IsGitlabEnterprise   bool             `json:"isGitLabEnterprise"`
	BitbucketClientKey   string           `json:"bitbucketClientKey,omitempty"`
	IsGithubEnterprise   bool             `json:"isGitHubEnterprise"`
	IsBitbucketServer    bool             `json:"isBitbucketServer"`
	FileName             string           `json:"fileName,omitempty"`
	IsTerragruntRunAll   bool             `json:"isTerragruntRunAll"`
	IsAzureDevOps        bool             `json:"isAzureDevOps" tfschema:"is_azure_devops"`
	IsHelmRepository     bool             `json:"isHelmRepository"`
	HelmChartName        string           `json:"helmChartName,omitempty"`
	TerragruntTfBinary   string           `json:"terragruntTfBinary,omitempty"`
	AnsibleVersion       string           `json:"ansibleVersion,omitempty"`
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

func (payload *TemplateCreatePayload) Invalidate() error {
	if payload.OrganizationId != "" {
		return errors.New("must not specify organizationId")
	}

	if payload.Type != TERRAGRUNT && payload.TerragruntVersion != "" {
		return errors.New("can't define terragrunt version for non-terragrunt template")
	}

	if payload.Type == TERRAGRUNT && payload.TerragruntVersion == "" {
		return errors.New("must supply terragrunt version")
	}

	if payload.Type == OPENTOFU && payload.OpentofuVersion == "" {
		return errors.New("must supply opentofu version")
	}

	if payload.TerragruntTfBinary != "" && payload.Type != TERRAGRUNT {
		return fmt.Errorf("terragrunt_tf_binary should only be used when the template type is 'terragrunt', but type is '%s'", payload.Type)
	}

	if payload.IsTerragruntRunAll {
		if payload.Type != TERRAGRUNT {
			return errors.New(`can't set is_terragrunt_run_all to "true" for non-terragrunt template`)
		}

		c, _ := semver.NewConstraint(">= 0.28.1")

		v, err := semver.NewVersion(payload.TerragruntVersion)
		if err != nil {
			return fmt.Errorf("invalid semver version %s: %w", payload.TerragruntVersion, err)
		}

		if !c.Check(v) {
			return errors.New("can't set is_terragrunt_run_all to 'true' for terragrunt versions lower than 0.28.1")
		}
	}

	if payload.Type == "ansible" && payload.AnsibleVersion != "latest" {
		if payload.AnsibleVersion == "" {
			return errors.New("'ansible_version' is required")
		}

		c, _ := semver.NewConstraint(">= 3.0.0")

		v, err := semver.NewVersion(payload.AnsibleVersion)
		if err != nil {
			return fmt.Errorf("invalid ansible version '%s': %w", payload.AnsibleVersion, err)
		}

		if !c.Check(v) {
			return errors.New("supported ansible versions are 3.0.0 and above")
		}
	}

	if payload.Type == "cloudformation" && payload.FileName == "" {
		return errors.New("file_name is required with cloudformation template type")
	}

	if payload.Type != "cloudformation" && payload.FileName != "" {
		return fmt.Errorf("file_name cannot be set when template type is: %s", payload.Type)
	}

	if payload.IsHelmRepository || payload.HelmChartName != "" {
		if payload.Type != "helm" {
			return errors.New(`can't set is_helm_repository to "true" for non-helm template`)
		}

		if payload.HelmChartName == "" {
			return errors.New("helm_chart_name is required with helm repository")
		}

		if !payload.IsHelmRepository {
			return errors.New(`is_helm_repositroy set to "true" is required with "helm_chart_name"`)
		}
	}

	if payload.Type != TERRAGRUNT && payload.Type != "terraform" {
		payload.TerraformVersion = ""
	}

	if payload.Type != OPENTOFU && payload.TerragruntTfBinary != OPENTOFU {
		payload.OpentofuVersion = ""
	}

	return nil
}

func (client *ApiClient) TemplateCreate(payload TemplateCreatePayload) (Template, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return Template{}, err
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
	return client.http.Delete("/blueprints/"+id, nil)
}

func (client *ApiClient) TemplateUpdate(id string, payload TemplateCreatePayload) (Template, error) {
	organizationId, err := client.OrganizationId()
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
	organizationId, err := client.OrganizationId()
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
	return client.http.Delete("/blueprints/"+templateId+"/projects/"+projectId, nil)
}

func (client *ApiClient) VariablesFromRepository(payload *VariablesFromRepositoryPayload) ([]ConfigurationVariable, error) {
	paramsInterface, err := toParamsInterface(payload)
	if err != nil {
		return nil, err
	}

	params := map[string]string{}

	for key, value := range paramsInterface {
		switch key {
		case "githubInstallationId":
			params[key] = strconv.Itoa(int(value.(float64)))
		case "sshKeyIds":
			sshkeys := []string{}

			if value != nil {
				for _, sshkey := range value.([]interface{}) {
					sshkeys = append(sshkeys, "\""+sshkey.(string)+"\"")
				}
			}

			params[key] = "[" + strings.Join(sshkeys, ",") + "]"
		default:
			params[key] = value.(string)
		}
	}

	var result []ConfigurationVariable
	if err := client.http.Get("/blueprints/variables-from-repository", params, &result); err != nil {
		return nil, err
	}

	return result, nil
}
