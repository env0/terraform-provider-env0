package env0

import (
	"context"
	"fmt"
	"strings"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTeamProjectAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamProjectAssignmentCreateOrUpdate,
		ReadContext:   resourceTeamProjectAssignmentRead,
		UpdateContext: resourceTeamProjectAssignmentCreateOrUpdate,
		DeleteContext: resourceTeamProjectAssignmentDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceTeamProjectAssignmentImport},

		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:        schema.TypeString,
				Description: "id of the team",
				Required:    true,
				ForceNew:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "id of the project",
				Required:    true,
				ForceNew:    true,
			},
			"role": {
				Type:             schema.TypeString,
				Description:      "the assigned built-in role (Admin, Planner, Viewer, Deployer)",
				Optional:         true,
				ValidateDiagFunc: NewRoleValidator([]string{"Admin", "Planner", "Viewer", "Deployer"}),
				ExactlyOneOf:     []string{"custom_role_id", "role"},
			},
			"custom_role_id": {
				Type:             schema.TypeString,
				Description:      "id of the assigned custom role",
				Optional:         true,
				ExactlyOneOf:     []string{"custom_role_id", "role"},
				ValidateDiagFunc: ValidateNotEmptyString,
			},
		},
	}
}

func resourceTeamProjectAssignmentCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.TeamRoleAssignmentCreateOrUpdatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	role, ok := d.GetOk("role")
	if !ok {
		role = d.Get("custom_role_id")
	}
	payload.Role = role.(string)

	assignment, err := apiClient.TeamRoleAssignmentCreateOrUpdate(&payload)
	if err != nil {
		return diag.Errorf("could not create assignment: %v", err)
	}

	d.SetId(assignment.Id)

	return nil
}

func resourceTeamProjectAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.TeamRoleAssignmentListPayload

	id := d.Id()
	payload.ProjectId = d.Get("project_id").(string)

	assignments, err := apiClient.TeamRoleAssignments(&payload)
	if err != nil {
		return diag.Errorf("could not get assignments: %v", err)
	}

	for _, assignment := range assignments {
		if assignment.Id == id {
			if err := writeResourceData(&assignment, d); err != nil {
				return diag.Errorf("schema resource data serialization failed: %v", err)
			}

			if client.IsBuiltinRole(assignment.Role) {
				d.Set("role", assignment.Role)
			} else {
				d.Set("custom_role_id", assignment.Role)
			}

			return nil
		}
	}

	tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]interface{}{"id": d.Id()})
	d.SetId("")

	return nil
}

func resourceTeamProjectAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.TeamRoleAssignmentDeletePayload

	payload.ProjectId = d.Get("project_id").(string)
	payload.TeamId = d.Get("team_id").(string)

	if err := apiClient.TeamRoleAssignmentDelete(&payload); err != nil {
		return diag.Errorf("could not delete assignment: %v", err)
	}

	return nil
}

func resourceTeamProjectAssignmentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiClient := meta.(client.ApiClientInterface)

	splitTeamProject := strings.Split(d.Id(), "_")
	if len(splitTeamProject) != 2 {
		return nil, fmt.Errorf("the id %v is invalid must be <team_id>_<project_id>", d.Id())
	}

	teamId := splitTeamProject[0]
	projectId := splitTeamProject[1]

	var payload client.TeamRoleAssignmentListPayload
	payload.ProjectId = projectId

	assignments, err := apiClient.TeamRoleAssignments(&payload)
	if err != nil {
		return nil, err
	}

	for _, assignment := range assignments {
		if assignment.TeamId == teamId {
			if err := writeResourceData(&assignment, d); err != nil {
				return nil, fmt.Errorf("schema resource data serialization failed: %w", err)
			}

			if client.IsBuiltinRole(assignment.Role) {
				d.Set("role", assignment.Role)
			} else {
				d.Set("custom_role_id", assignment.Role)
			}

			d.Set("project_id", projectId)

			return []*schema.ResourceData{d}, nil
		}
	}

	return nil, fmt.Errorf("assignment with id %v not found", d.Id())
}
