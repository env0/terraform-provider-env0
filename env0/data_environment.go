package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataEnvironment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataEnvironmentRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Description:  "the environment's id",
				Optional:     true,
				AtLeastOneOf: []string{"name", "id"},
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "the environments name",
				Optional:     true,
				AtLeastOneOf: []string{"name", "id"},
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "the environment's project id",
				Optional:    true,
			},
			"template_id": {
				Type:        schema.TypeString,
				Description: "the environment's template id",
				Computed:    true,
			},
		},
	}
}

func dataEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var err diag.Diagnostics
	var environment client.Environment

	id, ok := d.GetOk("id")
	if ok {
		environment, err = getEnvironmentById(id.(string), meta)
		if err != nil {
			return err
		}
	} else {
		name := d.Get("name")
		environment, err = getEnvironmentByName(name.(string), meta)
		if err != nil {
			return err
		}
	}

	d.SetId(environment.Id)
	setEnvironmentSchema(d, environment)
	return nil
}
