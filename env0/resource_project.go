package env0

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const PROJECT_DESTROY_TOTAL_WAIT_TIME = time.Minute * 10
const PROJECT_DESTROY_WAIT_INTERVAL = time.Second * 10

type ActiveEnvironmentError struct {
	message string
	retry   bool
}

func (e *ActiveEnvironmentError) Error() string {
	return e.message
}

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceProjectImport},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Description:      "name to give the project",
				Required:         true,
				ValidateDiagFunc: ValidateNotEmptyString,
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
				Description: "Destroy the project even when environments exist",
				Optional:    true,
				Default:     false,
			},
			"wait": {
				Type:        schema.TypeBool,
				Description: "Wait for all environments to be destroyed before destroying this project (up to 10 minutes)",
				Optional:    true,
				Default:     false,
			},
			"parent_project_id": {
				Type:        schema.TypeString,
				Description: "If set, the project becomes a 'sub-project' of the parent project. See https://docs.env0.com/docs/sub-projects",
				Optional:    true,
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.ProjectCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	project, err := apiClient.ProjectCreate(payload)
	if err != nil {
		return diag.Errorf("could not create project: %v", err)
	}

	d.SetId(project.Id)

	return nil
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	project, err := apiClient.Project(d.Id())
	if err != nil {
		return ResourceGetFailure(ctx, "project", d, err)
	}

	if err := writeResourceData(&project, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	return nil
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()

	var payload client.ProjectUpdatePayload

	if d.HasChange("parent_project_id") {
		parentProjectId := d.Get("parent_project_id").(string)

		if err := apiClient.ProjectMove(id, parentProjectId); err != nil {
			return diag.Errorf("could not move project: %v", err)
		}
	}

	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if _, err := apiClient.ProjectUpdate(id, payload); err != nil {
		return diag.Errorf("could not update project: %v", err)
	}

	return nil
}

func resourceProjectAssertCanDelete(d *schema.ResourceData, meta any) error {
	forceDestroy := d.Get("force_destroy").(bool)
	if forceDestroy {
		return nil
	}

	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()

	envs, err := apiClient.ProjectEnvironments(id)
	if err != nil {
		return err
	}

	for _, env := range envs {
		if env.IsArchived == nil || !*env.IsArchived {
			return &ActiveEnvironmentError{
				retry:   true,
				message: fmt.Sprintf("found an active environment %s (remove the environment or use the force_destroy flag)", env.Name),
			}
		}
	}

	return nil
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()

	if d.Get("wait").(bool) {
		waitInterval := PROJECT_DESTROY_WAIT_INTERVAL
		if os.Getenv("TF_ACC") == "1" { // For acceptance tests.
			waitInterval = time.Second
		}

		ticker := time.NewTicker(waitInterval)                  // When invoked check if project can be deleted.
		timer := time.NewTimer(PROJECT_DESTROY_TOTAL_WAIT_TIME) // When invoked wait time has elapsed.
		done := make(chan bool)

		go func() {
			for {
				select {
				case <-timer.C:
					done <- true

					return
				case <-ticker.C:
					err := resourceProjectAssertCanDelete(d, meta)
					if err != nil {
						if aeerr, ok := err.(*ActiveEnvironmentError); ok {
							if aeerr.retry {
								continue
							}
						}
					}

					done <- true

					return
				}
			}
		}()

		<-done
	}

	if err := resourceProjectAssertCanDelete(d, meta); err != nil {
		return diag.Errorf("could not delete project: %v", err)
	}

	if err := apiClient.ProjectDelete(id); err != nil {
		return diag.Errorf("could not delete project: %v", err)
	}

	return nil
}

func resourceProjectImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	id := d.Id()
	_, err := uuid.Parse(id)

	var project client.Project

	if err == nil {
		tflog.Info(ctx, "Resolving project by id", map[string]any{"id": id})

		if project, err = getProjectById(id, meta); err != nil {
			return nil, err
		}
	} else {
		tflog.Info(ctx, "Resolving project by name", map[string]any{"name": id})

		if project, err = getProjectByName(id, "", "", "", meta); err != nil {
			return nil, err
		}
	}

	if err := writeResourceData(&project, d); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
