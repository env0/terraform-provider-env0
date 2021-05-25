package env0

import (
	"context"
	"errors"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*client.ApiClient)

	name := d.Get("name").(string)
	project, err := apiClient.ProjectCreate(name)
	if err != nil {
		return diag.Errorf("could not create project: %v", err)
	}

	d.SetId(project.Id)

	return nil
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*client.ApiClient)

	id := d.Id()
	_, err := apiClient.Project(id)
	if err != nil {
		return diag.Errorf("could not get project: %v", err)
	}
	return nil
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("not implemented")
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*client.ApiClient)

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
		_, getErr = getProjectById(id, meta)
	} else {
		_, getErr = getProjectByName(id, meta)
	}
	if getErr != nil {
		return nil, errors.New(getErr[0].Summary)
	} else {
		return []*schema.ResourceData{d}, nil
	}
}
