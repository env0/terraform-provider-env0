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
			"include_archived_projects": {
				Type:        schema.TypeBool,
				Description: "set to 'true' to include archived projects (defaults to 'false')",
				Optional:    true,
				Default:     false,
			},
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
						"is_archived": {
							Type:        schema.TypeBool,
							Description: "'true' if the project is archived",
							Computed:    true,
						},
						"parent_project_id": {
							Type:        schema.TypeString,
							Description: "the parent project id (if one exist)",
							Computed:    true,
						},
						"hierarchy": {
							Type:        schema.TypeString,
							Description: "the project hierarchy (e.g. uuid1|uuid2|...)",
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

	includeArchivedProjects := d.Get("include_archived_projects").(bool)

	projects, err := apiClient.Projects()
	if err != nil {
		return diag.Errorf("failed to get list of projects: %v", err)
	}

	filteredProjects := []client.Project{}

	for _, project := range projects {
		if includeArchivedProjects || !project.IsArchived {
			filteredProjects = append(filteredProjects, project)
		}
	}

	if err := writeResourceDataSlice(filteredProjects, "projects", d); err != nil {
		return diag.Errorf("schema slice resource data serialization failed: %v", err)
	}

	d.SetId("projects")

	return nil
}
