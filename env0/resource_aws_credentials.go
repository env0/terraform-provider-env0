package env0

import (
	"context"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAwsCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAwsCredentialsCreate,
		ReadContext:   resourceAwsCredentialsRead,
		DeleteContext: resourceAwsCredentialsDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceAwsCredentialsImport},

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

	credentials, err := apiClient.AwsCredentialsCreate(request)
	if err != nil {
		return diag.Errorf("could not create credentials key: %v", err)
	}

	d.SetId(credentials.Id)

	return nil
}

func resourceAwsCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	credentials, err := apiClient.CloudCredentials(d.Id())
	if err != nil {
		return ResourceGetFailure("aws credentials", d, err)
	}

	d.Set("name", credentials.Name)
	d.SetId(credentials.Id)

	return nil
}

func resourceAwsCredentialsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	err := apiClient.CloudCredentialsDelete(id)
	if err != nil {
		return diag.Errorf("could not delete credentials: %v", err)
	}
	return nil
}

func resourceAwsCredentialsImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	credentials, err := getCredentials(d.Id(), "AWS_", meta)
	if err != nil {
		if _, ok := err.(*client.NotFoundError); ok {
			return nil, fmt.Errorf("aws credentials resource with id %v not found", d.Id())
		}
		return nil, err
	}

	d.Set("name", credentials.Name)
	d.SetId(credentials.Id)

	return []*schema.ResourceData{d}, nil
}
