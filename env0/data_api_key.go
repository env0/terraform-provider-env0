package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataApiKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataApiKeyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "the name of the api key",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "the id of the api key",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"api_key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the api key id",
			},
		},
	}
}

func dataApiKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var apiKey *client.ApiKey

	var err error

	id, ok := d.GetOk("id")
	if ok {
		apiKey, err = getApiKeyById(id.(string), meta)

		if err != nil {
			return diag.Errorf("could not read api key: %v", err)
		}
	} else {
		apiKey, err = getApiKeyByName(d.Get("name").(string), meta)
		if err != nil {
			return diag.Errorf("could not read api key: %v", err)
		}
	}

	if err := writeResourceData(apiKey, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}
