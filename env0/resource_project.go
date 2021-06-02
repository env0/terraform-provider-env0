package env0

import (
	"context"
	"errors"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceProjectImport},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name to give the project",
				Required:    true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "id of the project",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "description of the project",
				Optional:    true,
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	project, err := apiClient.ProjectCreate(name, description)
	if err != nil {
		return diag.Errorf("could not create project: %v", err)
	}

	d.SetId(project.Id)
	d.Set("name", project.Name)
	d.Set("description", project.Description)

	return nil
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	project, err := apiClient.Project(d.Id())
	if err != nil {
		return diag.Errorf("could not get project: %v", err)
	}

	d.Set("name", project.Name)
	d.Set("description", project.Description)

	return nil
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	payload := client.UpdateProjectPayload{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
	_, err := apiClient.ProjectUpdate(id, payload)
	if err != nil {
		return diag.Errorf("could not update project: %v", err)
	}
	return nil
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	err := apiClient.ProjectDelete(id)
	if err != nil {
		return diag.Errorf("could not delete project: %v", err)
	}
	return nil
}

func resourceProjectImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()
	var getErr diag.Diagnostics
	_, uuidErr := uuid.Parse(id)
	if uuidErr == nil {
		log.Println("[INFO] Resolving Project by id: ", id)
		_, getErr = getProjectById(id, meta)
	} else {
		log.Println("[DEBUG] ID is not a valid env0 id ", id)
		log.Println("[INFO] Resolving Project by name: ", id)
		var project client.Project
		project, getErr = getProjectByName(id, meta)
		d.SetId(project.Id)
	}
	if getErr != nil {
		return nil, errors.New(getErr[0].Summary)
	} else {
		return []*schema.ResourceData{d}, nil
	}
}
