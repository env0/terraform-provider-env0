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
	var credentials *client.Credentials

	id, ok := d.GetOk("id")
	if ok {
		credentials, err = getCostCredentialsById(id.(string), string(client.AzureCostCredentialsType), meta)
		if err != nil {
			return err
		}
	} else {
		name, _ := d.Get("name").(string) // name must be specified here
		credentials, err = getCostCredentialsByName(name, string(client.AzureCostCredentialsType), meta)
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
