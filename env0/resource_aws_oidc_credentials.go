package env0

import (
	"context"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAwsOidcCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAwsOidcCredentialsCreate,
		UpdateContext: resourceAwsOidcCredentialsUpdate,
		ReadContext:   resourceCredentialsRead(AWS_OIDC_TYPE),
		DeleteContext: resourceCredentialsDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceCredentialsImport(AWS_OIDC_TYPE)},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name for the oidc credentials",
				Required:    true,
				ForceNew:    true,
			},
			"role_arn": {
				Type:        schema.TypeString,
				Description: "the aws role arn",
				Required:    true,
			},
			"duration": {
				Type:             schema.TypeInt,
				Description:      "the session duration in seconds. If set must be one of the following: 3600 (1h), 7200 (2h), 14400 (4h), 18000 (5h default), 28800 (8h), 43200 (12h)",
				Optional:         true,
				ValidateDiagFunc: NewIntInValidator([]int{3600, 7200, 14400, 18000, 28800, 43200}),
				Default:          18000,
			},
		},
	}
}

func awsOidcCredentialsGetValue(d *schema.ResourceData) (client.AwsCredentialsValuePayload, error) {
	value := client.AwsCredentialsValuePayload{}

	if err := readResourceData(&value, d); err != nil {
		return value, fmt.Errorf("schema resource data deserialization failed: %w", err)
	}

	value.RoleArn = d.Get("role_arn").(string) // tfschema is set (for older resources) need to manually set the role arn.

	return value, nil
}

func resourceAwsOidcCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	value, err := awsOidcCredentialsGetValue(d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := client.AwsCredentialsCreatePayload{
		Name:  d.Get("name").(string),
		Value: value,
		Type:  client.AwsOidcCredentialsType,
	}

	credentials, err := apiClient.CredentialsCreate(&request)
	if err != nil {
		return diag.Errorf("could not create aws oidc credentials: %v", err)
	}

	d.SetId(credentials.Id)

	return nil
}

func resourceAwsOidcCredentialsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	value, err := awsOidcCredentialsGetValue(d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := client.AwsCredentialsCreatePayload{
		Value: value,
		Type:  client.AwsOidcCredentialsType,
	}

	if _, err := apiClient.CredentialsUpdate(d.Id(), &request); err != nil {
		return diag.Errorf("could not update aws oidc credentials: %s %v", d.Id(), err)
	}

	return nil
}
