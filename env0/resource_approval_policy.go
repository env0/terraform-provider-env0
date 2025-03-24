package env0

import (
	"context"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApprovalPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApprovalPolicyCreate,
		ReadContext:   resourceApprovalPolicyRead,
		UpdateContext: resourceApprovalPolicyUpdate,
		DeleteContext: resourceApprovalPolicyDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceApprovalPolicyImport},

		Schema: getConfigurationTemplateSchema(ApprovalPolicy),
	}
}

func resourceApprovalPolicyCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.ApprovalPolicyCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	approvalPolicy, err := apiClient.ApprovalPolicyCreate(&payload)
	if err != nil {
		return diag.Errorf("failed to create approval policy: %v", err)
	}

	d.SetId(approvalPolicy.Id)

	return nil
}

func resourceApprovalPolicyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	approvalPolicy, err := apiClient.Template(d.Id())
	if err != nil {
		return ResourceGetFailure(ctx, "approval policy", d, err)
	}

	if approvalPolicy.IsDeleted && !d.IsNewResource() {
		tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]any{"id": d.Id()})
		d.SetId("")

		return nil
	}

	if err := writeResourceData(&approvalPolicy, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceApprovalPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.ApprovalPolicyUpdatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if _, err := apiClient.ApprovalPolicyUpdate(&payload); err != nil {
		return diag.Errorf("failed to update approval policy: %v", err)
	}

	return nil
}

func resourceApprovalPolicyDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.TemplateDelete(d.Id()); err != nil {
		return diag.Errorf("could not delete approval policy template: %v", err)
	}

	return nil
}

func getApprovalPolicyByName(name string, meta any) (*client.ApprovalPolicy, error) {
	apiClient := meta.(client.ApiClientInterface)

	approvalPolicies, err := apiClient.ApprovalPolicies(name)
	if err != nil {
		return nil, err
	}

	if len(approvalPolicies) == 0 {
		return nil, fmt.Errorf("approval policy with name %v not found", name)
	}

	if len(approvalPolicies) > 1 {
		return nil, fmt.Errorf("found multiple approval policies with name: %s. Use id instead or make sure approval policy names are unique %v", name, approvalPolicies)
	}

	return &approvalPolicies[0], nil
}

func getApprovalPolicy(ctx context.Context, id string, meta any) (any, error) {
	if _, err := uuid.Parse(id); err == nil {
		tflog.Info(ctx, "Resolving approval policy by id", map[string]any{"id": id})

		template, err := meta.(client.ApiClientInterface).Template(id)
		if err != nil {
			return nil, err
		}

		if template.Type != string(ApprovalPolicy) {
			return nil, fmt.Errorf("template type requires type %s but received type %s", ApprovalPolicy, template.Type)
		}

		return &template, nil
	} else {
		tflog.Info(ctx, "Resolving approval policy by name", map[string]any{"name": id})

		return getApprovalPolicyByName(id, meta)
	}
}

func resourceApprovalPolicyImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	approvalPolicy, err := getApprovalPolicy(ctx, d.Id(), meta)
	if err != nil {
		return nil, err
	}

	if err := writeResourceData(approvalPolicy, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %w", err)
	}

	return []*schema.ResourceData{d}, nil
}
