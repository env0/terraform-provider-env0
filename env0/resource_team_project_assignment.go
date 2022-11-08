package env0

import (
	"context"
	"log"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTeamProjectAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamProjectAssignmentCreateOrUpdate,
		ReadContext:   resourceTeamProjectAssignmentRead,
		UpdateContext: resourceTeamProjectAssignmentCreateOrUpdate,
		DeleteContext: resourceTeamProjectAssignmentDelete,

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
				ValidateDiagFunc: ValidateRole,
				ExactlyOneOf:     []string{"custom_role_id", "role"},
			},
			"custom_role_id": {
				Type:         schema.TypeString,
				Description:  "id of the assigned custom role",
				Optional:     true,
				ExactlyOneOf: []string{"custom_role_id", "role"},
			},
		},
	}
}

func resourceTeamProjectAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	projectId := d.Get("project_id").(string)
	assignments, err := apiClient.TeamProjectAssignments(projectId)

	if err != nil {
		return diag.Errorf("could not get TeamProjectAssignment: %v", err)
	}

	found := false
	for _, assignment := range assignments {
		if assignment.Id == id {
			if err := writeResourceData(&assignment, d); err != nil {
				return diag.Errorf("schema resource data serialization failed: %v", err)
			}
			found = true
			break
		}
	}
	if !found && !d.IsNewResource() {
		log.Printf("[WARN] Drift Detected: Terraform will remove %s from state", d.Id())
		d.SetId("")
	}
	return nil
}

func resourceTeamProjectAssignmentCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.TeamProjectAssignmentPayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	response, err := apiClient.TeamProjectAssignmentCreateOrUpdate(payload)
	if err != nil {
		return diag.Errorf("could not Create or Update TeamProjectAssignment: %v", err)
	}

	d.SetId(response.Id)

	return nil
}

func resourceTeamProjectAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	err := apiClient.TeamProjectAssignmentDelete(d.Id())
	if err != nil {
		return diag.Errorf("could not delete TeamProjectAssignment: %v", err)
	}

	return nil
}
