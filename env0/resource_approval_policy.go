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

func resourceApprovalPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.TemplateCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	payload.Type = string(ApprovalPolicy)

	template, err := apiClient.TemplateCreate(payload)
	if err != nil {
		return diag.Errorf("could not create approval policy template: %v", err)
	}

	d.SetId(template.Id)

	return nil
}

func resourceApprovalPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	approvalPolicy, err := apiClient.Template(d.Id())
	if err != nil {
		return ResourceGetFailure(ctx, "approval policy", d, err)
	}

	if approvalPolicy.IsDeleted && !d.IsNewResource() {
		tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]interface{}{"id": d.Id()})
		d.SetId("")
		return nil
	}

	if err := writeResourceData(&approvalPolicy, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceApprovalPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	request, problem := templateCreatePayloadFromParameters("", d)
	if problem != nil {
		return problem
	}

	request.Type = string(ApprovalPolicy)

	_, err := apiClient.TemplateUpdate(d.Id(), request)
	if err != nil {
		return diag.Errorf("could not update approval policy template: %v", err)
	}

	return nil
}

func resourceApprovalPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.TemplateDelete(d.Id()); err != nil {
		return diag.Errorf("could not delete approval policy template: %v", err)
	}

	return nil
}

func getApprovalPolicyByName(name string, meta interface{}) (*client.ApprovalPolicy, error) {
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

func getApprovalPolicy(ctx context.Context, id string, meta interface{}) (interface{}, error) {
	if _, err := uuid.Parse(id); err == nil {
		tflog.Info(ctx, "Resolving approval policy by id", map[string]interface{}{"id": id})

		template, err := meta.(client.ApiClientInterface).Template(id)
		if err != nil {
			return nil, err
		}

		if template.Type != string(ApprovalPolicy) {
			return nil, fmt.Errorf("template type requires type %s but received type %s", ApprovalPolicy, template.Type)
		}

		return &template, nil
	} else {
		tflog.Info(ctx, "Resolving approval policy by name", map[string]interface{}{"name": id})

		return getApprovalPolicyByName(id, meta)
	}
}

func resourceApprovalPolicyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	approvalPolicy, err := getApprovalPolicy(ctx, d.Id(), meta)
	if err != nil {
		return nil, err
	}

	if err := writeResourceData(approvalPolicy, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %v", err)
	}

	return []*schema.ResourceData{d}, nil
}
