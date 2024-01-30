package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataProjectsRead,

		Schema: map[string]*schema.Schema{
			"projects": {
				Type:        schema.TypeList,
				Description: "list of projects",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "the name of the project",
							Computed:    true,
						},
						"id": {
							Type:        schema.TypeString,
							Description: "id of the project",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataProjectsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	projects, err := apiClient.Projects()
	if err != nil {
		return diag.Errorf("failed to get list of projects: %v", err)
	}

	if err := writeResourceDataSlice(projects, "projects", d); err != nil {
		return diag.Errorf("schema slice resource data serialization failed: %v", err)
	}

	d.SetId("projects")

	return nil
}
