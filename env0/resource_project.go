package env0

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/go-cty/cty"
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
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					name := i.(string)
					if name == "" {
						return diag.Errorf("Project name cannot be empty")
					}
					return nil
				},
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
			"force_destroy": {
				Type:        schema.TypeBool,
				Description: "A boolean that indicates if the project should be deleted if enviornments exist",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func setProjectSchema(d *schema.ResourceData, project client.Project) {
	d.Set("name", project.Name)
	d.Set("description", project.Description)
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	project, err := apiClient.ProjectCreate(client.ProjectCreatePayload{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	})
	if err != nil {
		return diag.Errorf("could not create project: %v", err)
	}

	d.SetId(project.Id)

	setProjectSchema(d, project)

	return nil
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	project, err := apiClient.Project(d.Id())
	if err != nil {
		return diag.Errorf("could not get project: %v", err)
	}

	setProjectSchema(d, project)

	return nil
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	payload := client.ProjectCreatePayload{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	project, err := apiClient.ProjectUpdate(id, payload)
	if err != nil {
		return diag.Errorf("could not update project: %v", err)
	}

	setProjectSchema(d, project)

	return nil
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()

	forceDestroy := d.Get("force_destroy").(bool)

	if !forceDestroy {
		// Wait up to a minute before returning an error.

		ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
		defer cancel()

		for {
			envs, err := apiClient.ProjectEnvironments(id)
			if err != nil {
				return diag.Errorf("could not get project environments: %v", err)
			}

			// Filter out archived (destroyed) environments.
			hasActiveEnvs := false

			for _, env := range envs {
				if !env.IsArchived {
					hasActiveEnvs = true
					break
				}
			}

			if !hasActiveEnvs {
				break
			}

			select {
			case <-ctx.Done():
				return diag.Errorf("unable to remove a project with environments (remove the environments or use the force_destroy flag)")
			case <-time.After(time.Second * 3):
			}
		}
	}

	if err := apiClient.ProjectDelete(id); err != nil {
		return diag.Errorf("could not delete project: %v", err)
	}
	return nil
}

func resourceProjectImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()
	var getErr diag.Diagnostics
	_, err := uuid.Parse(id)
	if err == nil {
		log.Println("[INFO] Resolving Project by id: ", id)
		_, getErr = getProjectById(id, meta)
	} else {
		log.Println("[DEBUG] ID is not a valid env0 id ", id)
		log.Println("[INFO] Resolving Project by name: ", id)

		var project client.Project
		project, getErr = getProjectByName(id, meta)

		d.SetId(project.Id)
		setProjectSchema(d, project)
	}

	if getErr != nil {
		return nil, errors.New(getErr[0].Summary)
	} else {
		return []*schema.ResourceData{d}, nil
	}
}
