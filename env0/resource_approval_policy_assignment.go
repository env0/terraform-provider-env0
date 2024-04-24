package env0

import (
	"context"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func setApprovalPolicyAssignmentId(d *schema.ResourceData, assignment *client.ApprovalPolicyAssignment) {
	d.SetId(fmt.Sprintf("%s|%s|%s", assignment.BlueprintId, assignment.Scope, assignment.ScopeId))
}

func resourceApprovalPolicyAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApprovalPolicyAssignmentCreate,
		ReadContext:   resourceApprovalPolicyAssignmentRead,
		DeleteContext: resourceApprovalPolicyAssignmentDelete,

		Schema: map[string]*schema.Schema{
			"scope": {
				Type:             schema.TypeString,
				Description:      "the type of the scope. Valid values: PROJECT or BLUEPRINT. Default value: PROJECT",
				Optional:         true,
				Default:          client.ApprovalPolicyProjectScope,
				ForceNew:         true,
				ValidateDiagFunc: NewStringInValidator([]string{"PROJECT", "BLUEPRINT"}),
			},
			"scope_id": {
				Type:        schema.TypeString,
				Description: "the id of the scope (E.g. project id or template id)",
				Required:    true,
				ForceNew:    true,
			},
			"blueprint_id": {
				Type:        schema.TypeString,
				Description: "the id of the approval policy",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceApprovalPolicyAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var assignment client.ApprovalPolicyAssignment
	if err := readResourceData(&assignment, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	template, err := apiClient.Template(assignment.BlueprintId)
	if err != nil {
		return diag.Errorf("unable to get template with id %s: %v", assignment.BlueprintId, err)
	}

	if template.Type != string(ApprovalPolicy) {
		return diag.Errorf("template with id %s is of type %s, but %s type is required", assignment.BlueprintId, template.Type, ApprovalPolicy)
	}

	if _, err := apiClient.ApprovalPolicyAssign(&assignment); err != nil {
		return diag.Errorf("could not assign approval policy %s to scope %s %s: %v", assignment.BlueprintId, assignment.Scope, assignment.ScopeId, err)
	}

	setApprovalPolicyAssignmentId(d, &assignment)

	return nil
}

func resourceApprovalPolicyAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var assignment client.ApprovalPolicyAssignment
	if err := readResourceData(&assignment, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	approvalPolicyByScopeArr, err := apiClient.ApprovalPolicyByScope(assignment.Scope, assignment.ScopeId)
	if err != nil {
		return ResourceGetFailure(ctx, "approval policy assignment", d, err)
	}

	found := false
	for _, approvalPolicyByScope := range approvalPolicyByScopeArr {
		if approvalPolicyByScope.ApprovalPolicy.Id == assignment.BlueprintId {
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]interface{}{"id": d.Id()})
		d.SetId("")
		return nil
	}

	return nil
}

func resourceApprovalPolicyAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	scope := d.Get("scope").(string)
	scopeId := d.Get("scope_id").(string)
	bluePrintId := d.Get("blueprint_id").(string)

	id := fmt.Sprintf("%s#%s#%s", scope, scopeId, bluePrintId)

	if err := apiClient.ApprovalPolicyUnassign(id); err != nil {
		return diag.Errorf("failed to unassign approval policy %s: %v", id, err)
	}

	return nil
}
