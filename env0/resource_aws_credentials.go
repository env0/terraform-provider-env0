package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAwsCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAwsCredentialsCreate,
		ReadContext:   resourceCredentialsRead(AWS_TYPE),
		DeleteContext: resourceCredentialsDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceCredentialsImport(AWS_TYPE)},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name for the credentials",
				Required:    true,
				ForceNew:    true,
			},
			"arn": {
				Type:          schema.TypeString,
				Description:   "the aws role arn",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"access_key_id"},
			},
			"access_key_id": {
				Type:          schema.TypeString,
				Description:   "the aws access key id",
				Optional:      true,
				Sensitive:     true,
				ForceNew:      true,
				ConflictsWith: []string{"arn", "duration"},
				RequiredWith:  []string{"secret_access_key"},
			},
			"secret_access_key": {
				Type:          schema.TypeString,
				Description:   "the aws access key secret. In case your organization is self-hosted, please use a secret reference in the shape of ${ssm:<secret-id>}",
				Optional:      true,
				Sensitive:     true,
				ForceNew:      true,
				ConflictsWith: []string{"arn"},
				RequiredWith:  []string{"access_key_id"},
			},
			"duration": {
				Type:             schema.TypeInt,
				Description:      "the session duration in seconds for AWS_ASSUMED_ROLE_FOR_DEPLOYMENT. If set must be one of the following: 3600 (1h), 7200 (2h), 14400 (4h), 18000 (5h default), 28800 (8h), 43200 (12h)",
				Optional:         true,
				ValidateDiagFunc: NewIntInValidator([]int{3600, 7200, 14400, 18000, 28800, 43200}),
				ForceNew:         true,
				ConflictsWith:    []string{"access_key_id"},
			},
		},
	}
}

func resourceAwsCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	_, accessKeyExist := d.GetOk("access_key_id")
	_, arnExist := d.GetOk("arn")

	if !accessKeyExist && !arnExist {
		// Due to "import" must be enforced here and not in the schema level.
		// This fields are only available during creation (will not be returned in read or import).
		return diag.Errorf("one of `access_key_id,arn` must be specified")
	}

	apiClient := meta.(client.ApiClientInterface)

	value := client.AwsCredentialsValuePayload{}
	if err := readResourceData(&value, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	requestType := client.AwsAssumedRoleCredentialsType
	if _, ok := d.GetOk("access_key_id"); ok {
		requestType = client.AwsAccessKeysCredentialsType
		// Duration applies only for "ASSUME ROLE"
		value.Duration = 0
	}

	request := client.AwsCredentialsCreatePayload{
		Name:  d.Get("name").(string),
		Value: value,
		Type:  requestType,
	}

	credentials, err := apiClient.CredentialsCreate(&request)
	if err != nil {
		return diag.Errorf("could not create credentials key: %v", err)
	}

	d.SetId(credentials.Id)

	return nil
}
