package env0tfprovider

import (
	"context"

	"github.com/env0/terraform-provider-env0/env0apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataProjectRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "the name of the project",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "id of the project",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"created_by": {
				Type:        schema.TypeString,
				Description: "textual description of the entity who created the project",
				Computed:    true,
			},
			"role": {
				Type:        schema.TypeString,
				Description: "role of the authenticated user (through api key) in the project",
				Computed:    true,
			},
		},
	}
}

func dataProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*env0apiclient.ApiClient)
	var err error
	var project env0apiclient.Project

	id, ok := d.GetOk("id")
	if ok {
		project, err = apiClient.Project(id.(string))
		if err != nil {
			return diag.Errorf("Could not query project by id: %v", err)
		}
	} else {
		name, ok := d.GetOk("name")
		if !ok {
			return diag.Errorf("Either 'name' or 'id' must be specified")
		}
		projects, err := apiClient.Projects()
		if err != nil {
			return diag.Errorf("Could not query project by name: %v", err)
		}
		for _, candidate := range projects {
			if candidate.Name == name.(string) {
				project = candidate
				break
			}
		}
		if project.Id == "" {
			return diag.Errorf("Could not find a project with name: %s", name)
		}
	}

	d.SetId(project.Id)
	d.Set("name", project.Name)
	d.Set("created_by", project.CreatedBy)
	d.Set("role", project.Role)

	return nil
}
