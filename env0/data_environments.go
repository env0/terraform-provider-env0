package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataEnvironments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataEnvironmentsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Description:  "project id to filter environments by",
				Optional:     true,
				ExactlyOneOf: []string{"organization_id"},
			},
			"organization_id": {
				Type:         schema.TypeString,
				Description:  "organization id to list all environments in the organization",
				Optional:     true,
				ExactlyOneOf: []string{"project_id"},
			},
			"include_archived_environments": {
				Type:        schema.TypeBool,
				Description: "set to 'true' to include archived environments (defaults to 'false')",
				Optional:    true,
				Default:     false,
			},
			"environments": {
				Type:        schema.TypeList,
				Description: "A list of environments",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "The ID of the environment",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the environment",
							Computed:    true,
						},
						"project_id": {
							Type:        schema.TypeString,
							Description: "The project ID the environment belongs to",
							Computed:    true,
						},
						// Add additional fields as necessary
					},
				},
			},
		},
	}
}

func dataEnvironmentsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	projectId := d.Get("project_id").(string)
	organizationId := d.Get("organization_id").(string)
	includeArchived := d.Get("include_archived_environments").(bool)

	var environments []client.Environment

	var err error

	if projectId != "" {
		environments, err = apiClient.ProjectEnvironments(projectId)
	} else {
		// organization_id must be set due to ExactlyOneOf
		environments, err = apiClient.OrganizationEnvironments(organizationId)
	}

	if err != nil {
		return diag.Errorf("failed to get list of environments: %v", err)
	}

	filtered := []client.Environment{}

	for _, env := range environments {
		if includeArchived {
			filtered = append(filtered, env)

			continue
		}

		if env.IsArchived == nil || !*env.IsArchived {
			filtered = append(filtered, env)
		}
	}

	if err := writeResourceDataSlice(filtered, "environments", d); err != nil {
		return diag.Errorf("schema slice resource data serialization failed: %v", err)
	}

	if projectId != "" {
		d.SetId("environments_" + projectId)
	} else {
		d.SetId("environments_org_" + organizationId)
	}

	return nil
}
