package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataCostCredentials(credType string) *schema.Resource {
	return &schema.Resource{
		ReadContext: dataCostCredentialsRead(credType),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "the name of the credential",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "the id of the credential",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
		},
	}
}

func dataCostCredentialsRead(credType string) func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		var err diag.Diagnostics
		var credentials *client.Credentials

		id, ok := d.GetOk("id")
		if ok {
			credentials, err = getCostCredentialsById(id.(string), credType, meta)
			if err != nil {
				return err
			}
		} else {
			name, ok := d.GetOk("name")
			if !ok {
				return diag.Errorf("Either 'name' or 'id' must be specified")
			}
			credentials, err = getCostCredentialsByName(name.(string), credType, meta)
			if err != nil {
				return err
			}
		}

		errorWhenWriteData := writeResourceData(credentials, d)
		if errorWhenWriteData != nil {
			return diag.Errorf("Error: %v", errorWhenWriteData)
		}

		return nil
	}
}

func getCostCredentialsByName(name interface{}, credType string, meta interface{}) (*client.Credentials, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)
	credentialsList, err := apiClient.CloudCredentialsList()
	if err != nil {
		return &client.Credentials{}, diag.Errorf("Could not query Cost Credentials by name: %v", err)
	}

	credentialsByNameAndType := make([]client.Credentials, 0)
	for _, candidate := range credentialsList {
		if candidate.Name == name.(string) && candidate.Type == credType {
			credentialsByNameAndType = append(credentialsByNameAndType, candidate)
		}
	}

	if len(credentialsByNameAndType) > 1 {
		return &client.Credentials{}, diag.Errorf("Found multiple Cost Credentials for name: %s", name)
	}
	if len(credentialsByNameAndType) == 0 {
		return &client.Credentials{}, diag.Errorf("Could not find Cost Credentials with name: %s", name)
	}
	return &credentialsByNameAndType[0], nil
}

func getCostCredentialsById(id string, credType string, meta interface{}) (*client.Credentials, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)
	credentials, err := apiClient.CloudCredentials(id)
	if credentials.Type != credType {
		return &client.Credentials{}, diag.Errorf("Found credentials which are not Cost Credentials: %v", credentials)
	}
	if err != nil {
		return &client.Credentials{}, diag.Errorf("Could not query Cost Credentials: %v", err)
	}
	return &credentials, nil
}
