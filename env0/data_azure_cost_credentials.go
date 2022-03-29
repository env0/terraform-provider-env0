package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataAzureCostCredentials() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataAzureCostCredentialsRead,

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

func dataAzureCostCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var err diag.Diagnostics
	var credentials client.Credentials

	id, ok := d.GetOk("id")
	if ok {
		credentials, err = getAzureCostCredentialsById(id.(string), meta)
		if err != nil {
			return err
		}
	} else {
		name, _ := d.Get("name").(string) // name must be specified here
		credentials, err = getAzureCostCredentialsByName(name, meta)
		if err != nil {
			return err
		}
	}

	d.SetId(credentials.Id)
	d.Set("name", credentials.Name)

	return nil
}

func getAzureCostCredentialsByName(name interface{}, meta interface{}) (client.Credentials, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)
	credentialsList, err := apiClient.CloudCredentialsList()
	if err != nil {
		return client.Credentials{}, diag.Errorf("Could not query Azure Credentials by name: %v", err)
	}

	credentialsByNameAndType := make([]client.Credentials, 0)
	for _, candidate := range credentialsList {
		if candidate.Name == name.(string) && isValidAzureCostCredentialsType(candidate.Type) {
			credentialsByNameAndType = append(credentialsByNameAndType, candidate)
		}
	}

	if len(credentialsByNameAndType) > 1 {
		return client.Credentials{}, diag.Errorf("Found multiple Azure cost Credentials for name: %s", name)
	}
	if len(credentialsByNameAndType) == 0 {
		return client.Credentials{}, diag.Errorf("Could not find Azure cost Credentials with name: %s", name)
	}
	return credentialsByNameAndType[0], nil
}

func getAzureCostCredentialsById(id string, meta interface{}) (client.Credentials, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)
	credentials, err := apiClient.CloudCredentials(id)
	if !isValidAzureCostCredentialsType(credentials.Type) {
		return client.Credentials{}, diag.Errorf("Found credentials which are not Azure Credentials: %v", credentials)
	}
	if err != nil {
		return client.Credentials{}, diag.Errorf("Could not query Azure Credentials: %v", err)
	}
	return credentials, nil
}

func isValidAzureCostCredentialsType(credentialsType string) bool {
	return client.AzureCredentialsType(credentialsType) == client.AzureCostCredentialsType
}
