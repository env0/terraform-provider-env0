package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataGoogleCostCredentials() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataGcpCredentialsRead,

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

func dataGoogleCostCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var err diag.Diagnostics
	var credentials client.ApiKey

	id, ok := d.GetOk("id")
	if ok {
		credentials, err = getGoogleCostCredentialsById(id.(string), meta)
		if err != nil {
			return err
		}
	} else {
		name, _ := d.Get("name").(string) // name must be specified here
		credentials, err = getGoogleCostCredentialsByName(name, meta)
		if err != nil {
			return err
		}
	}

	d.SetId(credentials.Id)
	d.Set("name", credentials.Name)

	return nil
}

func getGoogleCostCredentialsByName(name interface{}, meta interface{}) (client.ApiKey, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)
	credentialsList, err := apiClient.CloudCredentialsList()
	if err != nil {
		return client.ApiKey{}, diag.Errorf("Could not query GCP Credentials by name: %v", err)
	}

	credentialsByNameAndType := make([]client.ApiKey, 0)
	for _, candidate := range credentialsList {
		if candidate.Name == name.(string) && isValidGoogleCostCredentialsType(candidate.Type) {
			credentialsByNameAndType = append(credentialsByNameAndType, candidate)
		}
	}

	if len(credentialsByNameAndType) > 1 {
		return client.ApiKey{}, diag.Errorf("Found multiple Google cost Credentials for name: %s", name)
	}
	if len(credentialsByNameAndType) == 0 {
		return client.ApiKey{}, diag.Errorf("Could not find GCP Credentials with name: %s", name)
	}
	return credentialsByNameAndType[0], nil
}

func getGoogleCostCredentialsById(id string, meta interface{}) (client.ApiKey, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)
	credentials, err := apiClient.CloudCredentials(id)
	if !isValidGoogleCostCredentialsType(credentials.Type) {
		return client.ApiKey{}, diag.Errorf("Found credentials which are not GCP Credentials: %v", credentials)
	}
	if err != nil {
		return client.ApiKey{}, diag.Errorf("Could not query GCP Credentials: %v", err)
	}
	return credentials, nil
}

func isValidGoogleCostCredentialsType(credentialsType string) bool {
	return client.GcpCredentialsType(credentialsType) == client.GoogleCostCredentiassType
}
