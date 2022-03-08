package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
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
			"description": {
				Type:        schema.TypeString,
				Description: "textual description of the project",
				Computed:    true,
			},
		},
	}
}

func dataProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var err diag.Diagnostics
	var project client.Project

	id, ok := d.GetOk("id")
	if ok {
		project, err = getProjectById(id.(string), meta)
		if err != nil {
			return err
		}
	} else {
		name, ok := d.GetOk("name")
		if !ok {
			return diag.Errorf("Either 'name' or 'id' must be specified")
		}
		project, err = getProjectByName(name.(string), meta)
		if err != nil {
			return err
		}
	}

	d.SetId(project.Id)
	d.Set("name", project.Name)
	d.Set("created_by", project.CreatedBy)
	d.Set("role", project.Role)
	d.Set("description", project.Description)

	return nil
}

func getProjectByName(name interface{}, meta interface{}) (client.Project, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)
	projects, err := apiClient.Projects()
	if err != nil {
		return client.Project{}, diag.Errorf("Could not query project by name: %v", err)
	}

	projectsByName := make([]client.Project, 0)
	for _, candidate := range projects {
		if candidate.Name == name.(string) && !candidate.IsArchived {
			projectsByName = append(projectsByName, candidate)
		}
	}

	if len(projectsByName) > 1 {
		return client.Project{}, diag.Errorf("Found multiple Projects for name: %s. Use ID instead or make sure Project names are unique %v", name, projectsByName)
	}
	if len(projectsByName) == 0 {
		return client.Project{}, diag.Errorf("Could not find a project with name: %s", name)
	}
	return projectsByName[0], nil
}

func getProjectById(id string, meta interface{}) (client.Project, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)
	project, err := apiClient.Project(id)
	if err != nil {
		return client.Project{}, diag.Errorf("Could not query template: %v", err)
	}
	return project, nil
}
