package client

import "fmt"

type CloudAccountCreatePayload struct {
	Provider      string      `json:"provider"`
	Name          string      `json:"name"`
	Configuration interface{} `json:"configuration" tfschema:"-"`
}

type CloudAccountUpdatePayload struct {
	Name          string      `json:"name"`
	Configuration interface{} `json:"configuration" tfschema:"-"`
}

type AWSCloudAccountConfiguration struct {
	AccountId                   string   `json:"accountId"`
	BucketName                  string   `json:"bucketName"`
	Prefix                      string   `json:"prefix,omitempty"`
	Regions                     []string `json:"regions"`
	ShouldPrefixUnderLogsFolder bool     `json:"shouldPrefixUnderLogsFolder"`
}

type CloudAccount struct {
	Id            string      `json:"id"`
	Provider      string      `json:"provider"`
	Name          string      `json:"name"`
	Health        bool        `json:"health"`
	Configuration interface{} `json:"configuration" tfschema:"-"`
}

func (client *ApiClient) CloudAccountCreate(payload *CloudAccountCreatePayload) (*CloudAccount, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, fmt.Errorf("failed to get organization id: %w", err)
	}

	payloadWithOrganizationId := struct {
		*CloudAccountCreatePayload
		OrganizationId string `json:"organizationId"`
	}{
		payload,
		organizationId,
	}

	var cloudAccount CloudAccount
	if err := client.http.Post("/cloud/configurations", &payloadWithOrganizationId, &cloudAccount); err != nil {
		return nil, err
	}

	return &cloudAccount, nil
}

func (client *ApiClient) CloudAccountUpdate(id string, payload *CloudAccountUpdatePayload) (*CloudAccount, error) {
	var cloudAccount CloudAccount
	if err := client.http.Put("/cloud/configurations/"+id, payload, &cloudAccount); err != nil {
		return nil, err
	}

	return &cloudAccount, nil
}

func (client *ApiClient) CloudAccountDelete(id string) error {
	return client.http.Delete("/cloud/configurations/"+id, nil)
}

func (client *ApiClient) CloudAccount(id string) (*CloudAccount, error) {
	var cloudAccount CloudAccount

	if err := client.http.Get("/cloud/configurations/"+id, nil, &cloudAccount); err != nil {
		return nil, err
	}

	return &cloudAccount, nil
}

func (client *ApiClient) CloudAccounts() ([]CloudAccount, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, fmt.Errorf("failed to get organization id: %w", err)
	}

	var cloudAccounts []CloudAccount
	if err := client.http.Get("/cloud/configurations", map[string]string{"organizationId": organizationId}, &cloudAccounts); err != nil {
		return nil, err
	}

	return cloudAccounts, nil
}
