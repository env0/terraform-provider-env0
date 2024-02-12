package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataModuleTestingProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataModuleTestingProjectRead,
		Description: "Can be used to get the project_id for agent_project_assignment and cloud_credentials_project_assignment",

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "the module testing project id",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "the module testing project id",
				Computed:    true,
			},
		},
	}
}

func dataModuleTestingProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(client.ApiClientInterface)

	moduleTestingProject, err := client.ModuleTestingProject()
	if err != nil {
		return diag.Errorf("could not get module testing project: %v", err)
	}

	if err := writeResourceData(moduleTestingProject, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}
