package env0

import (
	"context"
	"regexp"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var externalIdRegex *regexp.Regexp = regexp.MustCompile(`^[A-Za-z0-9+=\-.,@:/]{2,1000}$`)

func IsExternalIdValid(d *schema.ResourceData, apiClient client.ApiClientInterface) (bool, error) {
	isSelfHosted, err := apiClient.IsOrganizationSelfHostedAgent()

	if err != nil {
		return false, err
	}

	if isSelfHosted {
		// Validating external id isn't required in self hosted organizations.
		return true, nil
	}

	externalId, ok := d.GetOk("external_id")
	if !ok || externalId.(string) == "" {
		// not in schema or empty - is valid.
		return true, nil
	}

	return externalIdRegex.MatchString(externalId.(string)), nil
}

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
			"external_id": {
				Type:          schema.TypeString,
				Description:   "the aws role external id",
				Optional:      true,
				Sensitive:     true,
				ForceNew:      true,
				ConflictsWith: []string{"access_key_id"},
				RequiredWith:  []string{"arn"},
			},
			"access_key_id": {
				Type:          schema.TypeString,
				Description:   "the aws access key id",
				Optional:      true,
				Sensitive:     true,
				ForceNew:      true,
				ConflictsWith: []string{"arn", "external_id"},
				RequiredWith:  []string{"secret_access_key"},
			},
			"secret_access_key": {
				Type:          schema.TypeString,
				Description:   "the aws access key secret",
				Optional:      true,
				Sensitive:     true,
				ForceNew:      true,
				ConflictsWith: []string{"arn", "external_id"},
				RequiredWith:  []string{"access_key_id"},
			},
		},
	}
}

func resourceAwsCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, accessKeyExist := d.GetOk("access_key_id")
	_, arnExist := d.GetOk("arn")
	if !accessKeyExist && !arnExist {
		// Due to "import" must be inforced here and not in the schema level.
		// This fields are only available during creation (will not be returned in read or import).
		return diag.Errorf("one of `access_key_id,arn` must be specified")
	}

	apiClient := meta.(client.ApiClientInterface)

	valid, err := IsExternalIdValid(d, apiClient)
	if err != nil {
		return diag.Errorf("failed to validate external_id: %v", err)
	}
	if !valid {
		return diag.Errorf("external_id is invalid")
	}

	value := client.AwsCredentialsValuePayload{}
	if err := readResourceData(&value, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	requestType := client.AwsAssumedRoleCredentialsType
	if _, ok := d.GetOk("access_key_id"); ok {
		requestType = client.AwsAccessKeysCredentialsType
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
