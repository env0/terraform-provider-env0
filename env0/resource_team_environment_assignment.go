package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTeamEnvironmentAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamEnvironmentAssignmentCreate,
		UpdateContext: resourceTeamEnvironmentAssignmentUpdate,
		ReadContext:   resourceTeamEnvironmentAssignmentRead,
		DeleteContext: resourceTeamEnvironmentAssignmentDelete,

		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:        schema.TypeString,
				Description: "id of the team",
				Required:    true,
				ForceNew:    true,
			},
			"environment_id": {
				Type:        schema.TypeString,
				Description: "id of the environment",
				Required:    true,
				ForceNew:    true,
			},
			"role_id": {
				Type:             schema.TypeString,
				Description:      "id of the assigned custom role. The following built-in roles can be passed as well: Viewer, Planner, Deployer, Admin",
				Required:         true,
				ValidateDiagFunc: ValidateNotEmptyString,
			},
		},
	}
}

func resourceTeamEnvironmentAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var newAssignment client.AssignTeamRoleToEnvironmentPayload
	if err := readResourceData(&newAssignment, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	client := meta.(client.ApiClientInterface)

	assignment, err := client.AssignTeamRoleToEnvironment(&newAssignment)
	if err != nil {
		return diag.Errorf("could not create assignment: %v", err)
	}

	d.SetId(assignment.Id)

	return nil
}

func resourceTeamEnvironmentAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(client.ApiClientInterface)

	id := d.Id()
	environmentId := d.Get("environment_id").(string)

	assignments, err := client.TeamRoleEnvironmentAssignments(environmentId)
	if err != nil {
		return diag.Errorf("could not get assignments: %v", err)
	}

	for _, assignment := range assignments {
		if assignment.Id == id {
			if err := writeResourceData(&assignment, d); err != nil {
				return diag.Errorf("schema resource data serialization failed: %v", err)
			}

			d.Set("role_id", assignment.Role)

			return nil
		}
	}

	tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]interface{}{"id": d.Id()})
	d.SetId("")

	return nil
}

func resourceTeamEnvironmentAssignmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var payload client.AssignTeamRoleToEnvironmentPayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	client := meta.(client.ApiClientInterface)

	assignment, err := client.AssignTeamRoleToEnvironment(&payload)
	if err != nil {
		return diag.Errorf("could not update assignment: %v", err)
	}

	d.SetId(assignment.Id)

	return nil
}

func resourceTeamEnvironmentAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	environmentId := d.Get("environment_id").(string)
	teamId := d.Get("team_id").(string)

	client := meta.(client.ApiClientInterface)

	if err := client.RemoveTeamRoleFromEnvironment(environmentId, teamId); err != nil {
		return diag.Errorf("could not delete assignment: %v", err)
	}

	return nil
}
