package env0

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataApprovalPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataApprovalPolicyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "The name of the approval policy",
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "ID of the approval policy",
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
		},
	}
}

func dataApprovalPolicyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if id, ok := d.GetOk("id"); ok {
		approvalPolicy, err := getApprovalPolicyById(id.(string), meta)
		if err != nil {
			return diag.Errorf("failed to get approval policy by id: %v", err)
		}

		return writeApprovalPolicyData(approvalPolicy, d)
	}

	approvalPolicy, err := getApprovalPolicyByName(d.Get("name").(string), meta)
	if err != nil {
		return diag.Errorf("failed to get approval policy by name: %v", err)
	}

	return writeApprovalPolicyData(approvalPolicy, d)
}

func writeApprovalPolicyData(approvalPolicy any, d *schema.ResourceData) diag.Diagnostics {
	if err := writeResourceData(approvalPolicy, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}
