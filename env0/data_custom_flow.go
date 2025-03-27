package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataCustomFlow() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataCustomFlowRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "The name of the custom flow",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "ID of the custom flow",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
		},
	}
}

func dataCustomFlowRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var err error

	var customFlow *client.CustomFlow

	id, ok := d.GetOk("id")
	if ok {
		customFlow, err = meta.(client.ApiClientInterface).CustomFlow(id.(string))
		if err != nil {
			return diag.Errorf("failed to get custom flow by id: %v", err)
		}
	} else {
		name := d.Get("name")

		customFlow, err = getCustomFlowByName(name.(string), meta)
		if err != nil {
			return diag.Errorf("failed to get custom flow by name: %v", err)
		}
	}

	if err := writeResourceData(customFlow, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}
