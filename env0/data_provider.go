package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataProvider() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataProviderRead,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Description:  "the type/name of the provider",
				Optional:     true,
				ExactlyOneOf: []string{"type", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "id of the provider",
				Optional:     true,
				ExactlyOneOf: []string{"type", "id"},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "the description of the provider",
				Computed:    true,
			},
		},
	}
}

func dataProviderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var provider *client.Provider

	var err error

	id, ok := d.GetOk("id")
	if ok {
		provider, err = meta.(client.ApiClientInterface).Provider(id.(string))
	} else {
		name := d.Get("type").(string)
		provider, err = getProviderByName(name, meta)
	}

	if err != nil {
		return DataGetFailure("provider", id, err)
	}

	if err := writeResourceData(provider, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}
