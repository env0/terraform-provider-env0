package client

import "strings"

type AwsCredentialsType string
type GcpCredentialsType string
type AzureCredentialsType string

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
	Name           string                       `json:"name"`
	OrganizationId string                       `json:"organizationId"`
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
	Name           string                     `json:"name"`
	OrganizationId string                     `json:"organizationId"`
	Type           AwsCredentialsType         `json:"type"`
	Value          AwsCredentialsValuePayload `json:"value"`
}

type AwsCredentialsValuePayload struct {
	RoleArn         string `json:"roleArn" tfschema:"arn"`
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type GoogleCostCredentialsCreatePayload struct {
	Name           string                            `json:"name"`
	OrganizationId string                            `json:"organizationId"`
	Type           GcpCredentialsType                `json:"type"`
	Value          GoogleCostCredentialsValuePayload `json:"value"`
}

type GoogleCostCredentialsValuePayload struct {
	TableId string `json:"tableid"`
	Secret  string `json:"secret"`
}

type GcpCredentialsCreatePayload struct {
	Name           string                     `json:"name"`
	OrganizationId string                     `json:"organizationId"`
	Type           GcpCredentialsType         `json:"type"`
	Value          GcpCredentialsValuePayload `json:"value"`
}

type GcpCredentialsValuePayload struct {
	ProjectId         string `json:"projectId"`
	ServiceAccountKey string `json:"serviceAccountKey"`
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

const (
	GoogleCostCredentialsType            GcpCredentialsType   = "GCP_CREDENTIALS"
	AzureCostCredentialsType             AzureCredentialsType = "AZURE_CREDENTIALS"
	AwsCostCredentialsType               AwsCredentialsType   = "AWS_ASSUMED_ROLE"
	AwsAssumedRoleCredentialsType        AwsCredentialsType   = "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT"
	AwsAccessKeysCredentialsType         AwsCredentialsType   = "AWS_ACCESS_KEYS_FOR_DEPLOYMENT"
	GcpServiceAccountCredentialsType     GcpCredentialsType   = "GCP_SERVICE_ACCOUNT_FOR_DEPLOYMENT"
	AzureServicePrincipalCredentialsType AzureCredentialsType = "AZURE_SERVICE_PRINCIPAL_FOR_DEPLOYMENT"
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
	err = client.http.Get("/credentials", map[string]string{"organizationId": organizationId}, &credentials)
	if err != nil {
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
	err = client.http.Post("/credentials", request, &result)
	if err != nil {
		return Credentials{}, err
	}
	return result, nil
}

func (client *ApiClient) CloudCredentialsDelete(id string) error {
	return client.http.Delete("/credentials/" + id)
}
