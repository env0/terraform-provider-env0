package env0

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const defaultProjectDestroyTimeout = time.Minute * 10
const projectDestroyWaitInterval = time.Second * 10

type ActiveEnvironmentError struct {
	message string
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
		CustomizeDiff: resourceProjectCustomizeDiff,

		Importer: &schema.ResourceImporter{StateContext: resourceProjectImport},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(defaultProjectDestroyTimeout),
		},

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
				Description: "Archive the project even when it still contains active environments. Note: env0 archives the project rather than deleting it, and it does not destroy the environments in it: every one of them is marked inactive and its continuous deployment, PR plans and scheduled deployments are disabled, while its cloud resources keep running (and billing). Archiving cannot be undone, and the call fails when the project has active sub-projects, with or without this flag",
				Optional:    true,
				Default:     false,
			},
			"wait": {
				Type:        schema.TypeBool,
				Description: "Wait for the project's environments to be destroyed or archived before deleting it. The wait is bounded by the 'delete' timeout of the 'timeouts' block (defaults to 10 minutes)",
				Optional:    true,
				Default:     false,
			},
			"parent_project_id": {
				Type:        schema.TypeString,
				Description: "If set, the project becomes a 'sub-project' of the parent project. Changing it moves the project, with every environment and sub-project under it, under the new parent in place. Inherited variables and role visibility change accordingly. Set to \"\" to make it a top-level project. See https://docs.env0.com/docs/sub-projects",
				Optional:    true,
			},
			"tags": {
				Type:        schema.TypeList,
				Description: "tags for the project",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceProjectCustomizeDiff(_ context.Context, d *schema.ResourceDiff, meta any) error {
	if d.Id() == "" || !d.HasChange("parent_project_id") {
		return nil
	}

	newParentId := d.Get("parent_project_id").(string)

	// Empty also means the new parent is not known yet at plan time.
	if newParentId == "" {
		return nil
	}

	if newParentId == d.Id() {
		return errors.New("a project cannot be its own parent")
	}

	apiClient := meta.(client.ApiClientInterface)

	newParent, err := apiClient.Project(newParentId)
	if err != nil {
		return fmt.Errorf("could not validate parent project '%s': %w", newParentId, err)
	}

	if slices.Contains(strings.Split(newParent.Hierarchy, "|"), d.Id()) {
		return fmt.Errorf("cannot move project under '%s': it is one of its own sub-projects", newParentId)
	}

	return nil
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
		if !isEnvironmentActive(env) {
			continue
		}

		return &ActiveEnvironmentError{
			message: fmt.Sprintf("found an active environment %s - its infrastructure may still exist (remove the environment or use the force_destroy flag)", env.Name),
		}
	}

	return nil
}

// An active environment is one whose infrastructure may still exist: not archived, and not left
// INACTIVE without being archived by a destroy that succeeded.
func isEnvironmentActive(env client.Environment) bool {
	if env.IsArchived != nil && *env.IsArchived {
		return false
	}

	return env.Status != "INACTIVE"
}

// orphanedEnvironmentsWarning names the active environments that archiving the project will leave
// running with no env0 automation. The plan says only that the project will be destroyed, so without
// this the user never learns their infrastructure survived. Call it before the delete: afterwards
// every environment is archived and the list has nothing to report.
func orphanedEnvironmentsWarning(d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	envs, err := apiClient.ProjectEnvironments(d.Id())
	if err != nil {
		// force_destroy means "archive it regardless", so a list that fails downgrades the warning
		// rather than failing a delete that would otherwise go through.
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "could not list the environments this project will orphan",
			Detail:   fmt.Sprintf("Archiving project '%s' leaves any active environment in it running with no env0 automation, and listing them failed: %v. Check the project in env0.", d.Get("name").(string), err),
		}}
	}

	var names []string

	for _, env := range envs {
		if isEnvironmentActive(env) {
			names = append(names, env.Name)
		}
	}

	if len(names) == 0 {
		return nil
	}

	return diag.Diagnostics{{
		Severity: diag.Warning,
		Summary:  "env0 did not destroy the environments in this project",
		Detail: fmt.Sprintf("Archiving project '%s' marked these active environments inactive and disabled their continuous deployment, PR plans and scheduled deployments, but their cloud resources keep running (and billing): %s. That list was read just before the archive, so check the project in env0 for the full set and destroy them there to remove the infrastructure.",
			d.Get("name").(string), strings.Join(names, ", ")),
	}}
}

// waitForProjectEnvironmentsToBeArchived polls until no environment blocks the project delete, the
// timeout elapses or ctx is cancelled. A timeout names the environment that was still blocking.
func waitForProjectEnvironmentsToBeArchived(ctx context.Context, d *schema.ResourceData, meta any, timeout time.Duration) error {
	waitInterval := projectDestroyWaitInterval

	if os.Getenv("TF_ACC") == "1" { // For acceptance tests reducing interval to 1 second and clamping timeout to 10 seconds.
		waitInterval = time.Second
		timeout = min(timeout, time.Second*10)
	}

	ticker := time.NewTicker(waitInterval) // When invoked - check whether the project can be deleted.
	defer ticker.Stop()

	timer := time.NewTimer(timeout) // When invoked - timeout.
	defer timer.Stop()

	for {
		err := resourceProjectAssertCanDelete(d, meta)
		if err == nil {
			return nil
		}

		var activeEnvironmentError *ActiveEnvironmentError
		if !errors.As(err, &activeEnvironmentError) {
			return err
		}

		tflog.Info(ctx, "waiting for the project's environments to be destroyed or archived", map[string]any{"projectId": d.Id(), "reason": err.Error()})

		select {
		case <-timer.C:
			return fmt.Errorf("timeout! %w", err)
		case <-ctx.Done():
			// Under a real apply this is the branch that fires: the SDK wraps Delete in
			// context.WithTimeout(ctx, d.Timeout(TimeoutDelete)), so the same deadline is on ctx and starts
			// fractionally earlier. The timer only bounds direct callers that pass a context without one.
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return fmt.Errorf("timeout! %w", err)
			}

			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()

	var orphanWarning diag.Diagnostics

	// force_destroy short-circuits the assert, so there is nothing for the wait to do either.
	if d.Get("force_destroy").(bool) {
		orphanWarning = orphanedEnvironmentsWarning(d, meta)
	} else if d.Get("wait").(bool) {
		if err := waitForProjectEnvironmentsToBeArchived(ctx, d, meta, d.Timeout(schema.TimeoutDelete)); err != nil {
			return diag.Errorf("could not delete project: %v", err)
		}
	} else if err := resourceProjectAssertCanDelete(d, meta); err != nil {
		return diag.Errorf("could not delete project: %v", err)
	}

	// Errors here reach the user as the server wrote them, "active sub-projects" included.
	if err := apiClient.ProjectDelete(id); err != nil {
		return diag.Errorf("could not delete project: %v", err)
	}

	return orphanWarning
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
