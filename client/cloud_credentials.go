package client

import "strings"

type AwsCredentialsType string
type GcpCredentialsType string
type AzureCredentialsType string
type VaultCredentialsType string

type Credentials struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	OrganizationId string `json:"organizationId"`
	Type           string `json:"type"`
}

func (c *Credentials) HasPrefix(prefixList []string) bool {
	for _, prefix := range prefixList {
		if strings.HasPrefix(c.Type, prefix) {
			return true
		}
	}

	return false
}

type CredentialCreatePayload interface {
	SetOrganizationId(organizationId string)
}

type AzureCredentialsCreatePayload struct {
	Name           string                       `json:"name,omitempty"`
	OrganizationId string                       `json:"organizationId,omitempty"`
	Type           AzureCredentialsType         `json:"type"`
	Value          AzureCredentialsValuePayload `json:"value"`
}

type AzureCredentialsValuePayload struct {
	ClientId       string `json:"clientId"`
	ClientSecret   string `json:"clientSecret"`
	SubscriptionId string `json:"subscriptionId"`
	TenantId       string `json:"tenantId"`
}

type AwsCredentialsCreatePayload struct {
	Name           string                     `json:"name,omitempty"`
	OrganizationId string                     `json:"organizationId,omitempty"`
	Type           AwsCredentialsType         `json:"type"`
	Value          AwsCredentialsValuePayload `json:"value"`
}

type AwsCredentialsValuePayload struct {
	RoleArn         string `json:"roleArn" tfschema:"arn"`
	Duration        int    `json:"duration,omitempty"`
	AccessKeyId     string `json:"accessKeyId,omitempty"`
	SecretAccessKey string `json:"secretAccessKey,omitempty"`
}

type GoogleCostCredentialsCreatePayload struct {
	Name           string                            `json:"name,omitempty"`
	OrganizationId string                            `json:"organizationId,omitempty"`
	Type           GcpCredentialsType                `json:"type"`
	Value          GoogleCostCredentialsValuePayload `json:"value"`
}

type GoogleCostCredentialsValuePayload struct {
	TableId string `json:"tableid"`
	Secret  string `json:"secret"`
}

type GcpCredentialsCreatePayload struct {
	Name           string                     `json:"name,omitempty"`
	OrganizationId string                     `json:"organizationId,omitempty"`
	Type           GcpCredentialsType         `json:"type"`
	Value          GcpCredentialsValuePayload `json:"value"`
}

type GcpCredentialsValuePayload struct {
	ProjectId                          string `json:"projectId,omitempty"`
	ServiceAccountKey                  string `json:"serviceAccountKey,omitempty"`
	CredentialConfigurationFileContent string `json:"credentialConfigurationFileContent,omitempty"`
}

type VaultCredentialsValuePayload struct {
	Address            string `json:"address"`
	JwtAuthBackendPath string `json:"jwtAuthBackendPath"`
	RoleName           string `json:"roleName"`
	Version            string `json:"version"`
	Namespace          string `json:"namespace,omitempty"`
}

type VaultCredentialsCreatePayload struct {
	Name           string                       `json:"name,omitempty"`
	OrganizationId string                       `json:"organizationId,omitempty"`
	Type           VaultCredentialsType         `json:"type"`
	Value          VaultCredentialsValuePayload `json:"value"`
}

func (c *GoogleCostCredentialsCreatePayload) SetOrganizationId(organizationId string) {
	c.OrganizationId = organizationId
}

func (c *AwsCredentialsCreatePayload) SetOrganizationId(organizationId string) {
	c.OrganizationId = organizationId
}

func (c *GcpCredentialsCreatePayload) SetOrganizationId(organizationId string) {
	c.OrganizationId = organizationId
}

func (c *AzureCredentialsCreatePayload) SetOrganizationId(organizationId string) {
	c.OrganizationId = organizationId
}

func (c *VaultCredentialsCreatePayload) SetOrganizationId(organizationId string) {
	c.OrganizationId = organizationId
}

const (
	AwsCostCredentialsType               AwsCredentialsType   = "AWS_ASSUMED_ROLE"
	AwsAssumedRoleCredentialsType        AwsCredentialsType   = "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT"
	AwsAccessKeysCredentialsType         AwsCredentialsType   = "AWS_ACCESS_KEYS_FOR_DEPLOYMENT"
	AwsOidcCredentialsType               AwsCredentialsType   = "AWS_OIDC"
	GoogleCostCredentialsType            GcpCredentialsType   = "GCP_CREDENTIALS"
	GcpServiceAccountCredentialsType     GcpCredentialsType   = "GCP_SERVICE_ACCOUNT_FOR_DEPLOYMENT"
	GcpOidcCredentialsType               GcpCredentialsType   = "GCP_OIDC"
	AzureCostCredentialsType             AzureCredentialsType = "AZURE_CREDENTIALS"
	AzureServicePrincipalCredentialsType AzureCredentialsType = "AZURE_SERVICE_PRINCIPAL_FOR_DEPLOYMENT"
	AzureOidcCredentialsType             AzureCredentialsType = "AZURE_OIDC"
	VaultOidcCredentialsType             VaultCredentialsType = "VAULT_OIDC"
)

func (client *ApiClient) CloudCredentials(id string) (Credentials, error) {
	var credentials, err = client.CloudCredentialsList()
	if err != nil {
		return Credentials{}, err
	}

	for _, v := range credentials {
		if v.Id == id {
			return v, nil
		}
	}

	return Credentials{}, &NotFoundError{}
}

func (client *ApiClient) CloudCredentialsList() ([]Credentials, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return []Credentials{}, err
	}

	var credentials []Credentials

	if err := client.http.Get("/credentials", map[string]string{"organizationId": organizationId}, &credentials); err != nil {
		return []Credentials{}, err
	}

	return credentials, nil
}

func (client *ApiClient) CredentialsCreate(request CredentialCreatePayload) (Credentials, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return Credentials{}, err
	}

	request.SetOrganizationId(organizationId)

	var result Credentials
	if err := client.http.Post("/credentials", request, &result); err != nil {
		return Credentials{}, err
	}

	return result, nil
}

func (client *ApiClient) CredentialsUpdate(id string, request CredentialCreatePayload) (Credentials, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return Credentials{}, err
	}

	request.SetOrganizationId(organizationId)

	var result Credentials

	if err := client.http.Patch("/credentials/"+id, request, &result); err != nil {
		return Credentials{}, err
	}

	return result, nil
}

func (client *ApiClient) CloudCredentialsDelete(id string) error {
	return client.http.Delete("/credentials/"+id, nil)
}
