package env0

import (
	"context"
	"math"
	"time"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOrganizationPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationPolicyCreateOrUpdate,
		ReadContext:   resourceOrganizationPolicyRead,
		UpdateContext: resourceOrganizationPolicyCreateOrUpdate,
		DeleteContext: resourceOrganizationPolicyDelete,

		Schema: map[string]*schema.Schema{
			"max_ttl": {
				Type:             schema.TypeString,
				Description:      "the maximum environment time-to-live allowed on deploy time. Format is <number>-<M/w/d/h> (Examples: 12-h, 3-d, 1-w, 1-M). Omit for infinite ttl. must be equal or longer than default_ttl",
				Optional:         true,
				ValidateDiagFunc: ValidateTtl,
			},
			"default_ttl": {
				Type:             schema.TypeString,
				Description:      "the default environment time-to-live allowed on deploy time. Format is <number>-<M/w/d/h> (Examples: 12-h, 3-d, 1-w, 1-M). Omit for infinite ttl. must be equal or shorter than max_ttl",
				Optional:         true,
				ValidateDiagFunc: ValidateTtl,
			},
			"do_not_report_skipped_status_checks": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"do_not_consider_merge_commits_for_pr_plans": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enable_oidc": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "set to 'true' to enable OIDC token (JWT) availability during env0 deployments (defaults to 'false')",
			},
		},
	}
}

func resourceOrganizationPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	organization, err := apiClient.Organization()
	if err != nil {
		return diag.Errorf("could not get organization (for organization policy): %v", err)
	}

	if err := writeResourceData(&organization, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func getOrganizationPolicyTtl(value *string) (time.Duration, error) {
	if value == nil || *value == "" {
		return math.MaxInt64, nil
	}

	return ttlToDuration(*value)
}

func resourceOrganizationPolicyCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.OrganizationPolicyUpdatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	// Validate that default ttl is "less than or equal" max ttl.
	defaultTtl, err := getOrganizationPolicyTtl(payload.DefaultTtl)
	if err != nil {
		return diag.Errorf("invalid default ttl: %v", err)
	}

	maxTtl, err := getOrganizationPolicyTtl(payload.MaxTtl)
	if err != nil {
		return diag.Errorf("invalid max ttl: %v", err)
	}

	if maxTtl < defaultTtl {
		return diag.Errorf("default ttl must not be larger than max ttl: %d %d", defaultTtl, maxTtl)
	}

	organization, err := apiClient.OrganizationPolicyUpdate((payload))
	if err != nil {
		return diag.Errorf("could not update organization policy: %v", err)
	}

	d.SetId(organization.Id)

	return nil
}

func resourceOrganizationPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	// In cases of a "DELETE", update the organization policy to default values.
	var payload client.OrganizationPolicyUpdatePayload
	if _, err := apiClient.OrganizationPolicyUpdate(payload); err != nil {
		return diag.Errorf("could not update organization policy to default values: %v", err)
	}

	return nil
}
