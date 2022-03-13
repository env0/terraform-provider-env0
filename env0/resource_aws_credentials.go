package env0

import (
	"context"
	"errors"
	"log"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
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
				ExactlyOneOf:  []string{"access_key_id"},
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
				ExactlyOneOf:  []string{"arn"},
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
	apiClient := meta.(client.ApiClientInterface)

	value := client.AwsCredentialsValuePayload{}
	requestType := client.AwsAssumedRoleCredentialsType
	if arn, ok := d.GetOk("arn"); ok {
		value.RoleArn = arn.(string)
		requestType = client.AwsAssumedRoleCredentialsType
	}
	if externalId, ok := d.GetOk("external_id"); ok {
		value.ExternalId = externalId.(string)
	}
	if accessKeyId, ok := d.GetOk("access_key_id"); ok {
		value.AccessKeyId = accessKeyId.(string)
		requestType = client.AwsAccessKeysCredentialsType
	}
	if secretAccessKey, ok := d.GetOk("secret_access_key"); ok {
		value.SecretAccessKey = secretAccessKey.(string)
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

	id := d.Id()
	_, err := apiClient.CloudCredentials(id)
	if err != nil {
		return diag.Errorf("could not get credentials: %v", err)
	}
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
	id := d.Id()
	var getErr diag.Diagnostics
	_, uuidErr := uuid.Parse(id)
	if uuidErr == nil {
		log.Println("[INFO] Resolving AWS Credentials by id: ", id)
		_, getErr = getAwsCredentialsById(id, meta)
	} else {
		log.Println("[DEBUG] ID is not a valid env0 id ", id)
		log.Println("[INFO] Resolving AWS Credentials by name: ", id)
		var awsCredential client.ApiKey
		awsCredential, getErr = getAwsCredentialsByName(id, meta)
		d.SetId(awsCredential.Id)
	}
	if getErr != nil {
		return nil, errors.New(getErr[0].Summary)
	} else {
		return []*schema.ResourceData{d}, nil
	}
}
