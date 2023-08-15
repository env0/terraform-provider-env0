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
				Type:        schema.TypeString,
				Description: "the type of the scope. Valid values: PROJECT. Default value: PROJECT",
				Optional:    true,
				Default:     client.ApprovalPolicyProjectScope,
				ForceNew:    true,
			},
			"scope_id": {
				Type:        schema.TypeString,
				Description: "the id of the scope (E.g. project id)",
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

	scope := d.Get("scope").(string)
	scopeId := d.Get("scope_id").(string)
	blueprintId := d.Get("blueprint_id").(string)

	template, err := apiClient.Template(blueprintId)
	if err != nil {
		return diag.Errorf("unable to get template with id %s: %v", blueprintId, err)
	}

	if template.Type != string(ApprovalPolicy) {
		return diag.Errorf("template with id %s is of type %s, but %s type is required", blueprintId, template.Type, ApprovalPolicy)
	}

	assignment := client.ApprovalPolicyAssignment{
		Scope:       client.ApprovalPolicyAssignmentScope(scope),
		ScopeId:     scopeId,
		BlueprintId: blueprintId,
	}

	if _, err := apiClient.ApprovalPolicyAssign(&assignment); err != nil {
		return diag.Errorf("could not assign approval policy %s to scope %s %s: %v", blueprintId, scope, scopeId, err)
	}

	setApprovalPolicyAssignmentId(d, &assignment)

	return nil
}

func resourceApprovalPolicyAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	scope := d.Get("scope").(string)
	scopeId := d.Get("scope_id").(string)
	blueprintId := d.Get("blueprint_id").(string)

	approvalPolicyByScopeArr, err := apiClient.ApprovalPolicyByScope(scope, scopeId)
	if err != nil {
		return ResourceGetFailure(ctx, "approval policy assignment", d, err)
	}

	found := false
	for _, approvalPolicyByScope := range approvalPolicyByScopeArr {
		if approvalPolicyByScope.ApprovalPolicy.Id == blueprintId {
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]interface{}{"id": d.Id()})
		d.SetId("")
		return nil
	}

	assignment := client.ApprovalPolicyAssignment{
		Scope:       client.ApprovalPolicyAssignmentScope(scope),
		ScopeId:     scopeId,
		BlueprintId: blueprintId,
	}

	setApprovalPolicyAssignmentId(d, &assignment)

	return nil
}

func resourceApprovalPolicyAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	scope := d.Get("scope").(string)
	scopeId := d.Get("scope_id").(string)

	if err := apiClient.ApprovalPolicyUnassign(scope, scopeId); err != nil {
		return diag.Errorf("failed to unassign approval policy from scope %s %s: %v", scope, scopeId, err)
	}

	return nil
}
