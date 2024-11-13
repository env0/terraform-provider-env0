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

func setCustomFlowAssignmentId(d *schema.ResourceData, assignment client.CustomFlowAssignment) {
	d.SetId(fmt.Sprintf("%s|%s", assignment.ScopeId, assignment.BlueprintId))
}

func getCustomFlowAssignmentFromId(d *schema.ResourceData) (client.CustomFlowAssignment, error) {
	id := d.Id()
	split := strings.Split(id, "|")

	var assignment client.CustomFlowAssignment

	assignment.ScopeId = split[0]

	if len(split) > 1 {
		assignment.BlueprintId = split[1]
	} else {
		return assignment, fmt.Errorf("invalid id %s", id)
	}

	assignment.Scope = client.CustomFlowProjectScope

	return assignment, nil
}

func resourceCustomFlowAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomFlowAssignmentCreate,
		ReadContext:   resourceCustomFlowAssignmentRead,
		DeleteContext: resourceCustomFlowAssignmentDelete,

		Schema: map[string]*schema.Schema{
			"scope": {
				Type: schema.TypeString,
				// Note: at the moment the only valid scope is "PROJECT". May add more scopes in the future.
				Description: "the type of the scope. Valid values: PROJECT. Default value: PROJECT",
				Optional:    true,
				Default:     client.CustomFlowProjectScope,
				ForceNew:    true,
			},
			"scope_id": {
				Type:        schema.TypeString,
				Description: "the id of the scope (E.g. project id)",
				Required:    true,
				ForceNew:    true,
			},
			"template_id": {
				Type:        schema.TypeString,
				Description: "the id of the custom flow",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceCustomFlowAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var assignment client.CustomFlowAssignment
	if err := readResourceData(&assignment, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if err := apiClient.CustomFlowAssign([]client.CustomFlowAssignment{assignment}); err != nil {
		return diag.Errorf("could not assign custom flow to %s: %v", strings.ToLower(string(assignment.Scope)), err)
	}

	setCustomFlowAssignmentId(d, assignment)

	return nil
}

func resourceCustomFlowAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	assignmentFromId, err := getCustomFlowAssignmentFromId(d)
	if err != nil {
		return diag.Errorf("could get assignment from id: %v", err)
	}

	assignments, err := apiClient.CustomFlowGetAssignments([]client.CustomFlowAssignment{assignmentFromId})
	if err != nil {
		return diag.Errorf("could not get custom flow assignments for id %s: %v", assignmentFromId.ScopeId, err)
	}

	found := false

	for _, assignment := range assignments {
		if assignment.BlueprintId == assignmentFromId.BlueprintId {
			found = true

			break
		}
	}

	if !found && !d.IsNewResource() {
		tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]interface{}{"id": d.Id()})
		d.SetId("")

		return nil
	}

	setCustomFlowAssignmentId(d, assignmentFromId)

	return nil
}

func resourceCustomFlowAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	assignmentFromId, err := getCustomFlowAssignmentFromId(d)
	if err != nil {
		return diag.Errorf("could get assignment from id: %v", err)
	}

	if err := apiClient.CustomFlowUnassign([]client.CustomFlowAssignment{assignmentFromId}); err != nil {
		return diag.Errorf("failed to unassign %s from custom flow %s: %v", assignmentFromId.ScopeId, assignmentFromId.BlueprintId, err)
	}

	return nil
}
