package env0

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAwsCostCredentials() *schema.Resource {
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
				Type:        schema.TypeString,
				Description: "the aws role arn",
				Required:    true,
				ForceNew:    true,
			},
			"external_id": {
				Type:        schema.TypeString,
				Description: "the aws role external id",
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
			},
		},
	}
}

func resourceAwsCostCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)
	request := client.AwsCredentialsCreatePayload{
		Name: d.Get("name").(string),
		Value: client.AwsCredentialsValuePayload{
			RoleArn:    d.Get("arn").(string),
			ExternalId: d.Get("external_id").(string),
		},
		Type: client.AwsCostCredentialsType,
	}
	credentials, err := apiClient.AwsCredentialsCreate(request)
	if err != nil {
		return diag.Errorf("could not create credentials key: %v", err)
	}

	d.SetId(credentials.Id)

	return nil
}

func resourceAwsCostCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	_, err := apiClient.CloudCredentials(id)
	if err != nil {
		return diag.Errorf("could not get credentials: %v", err)
	}
	return nil
}

func resourceAwsCostCredentialsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	err := apiClient.CloudCredentialsDelete(id)
	if err != nil {
		return diag.Errorf("could not delete credentials: %v", err)
	}
	return nil
}

func resourceAwsCostCredentialsImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()
	var getErr diag.Diagnostics
	_, uuidErr := uuid.Parse(id)
	if uuidErr == nil {
		log.Println("[INFO] Resolving AWS Credentials by id: ", id)
		_, getErr = getAwsCostCredentialsById(id, meta)
	} else {
		log.Println("[DEBUG] ID is not a valid env0 id ", id)
		log.Println("[INFO] Resolving AWS Credentials by name: ", id)
		var awsCredential client.ApiKey
		awsCredential, getErr = getAwsCostCredentialsByName(id, meta)
		d.SetId(awsCredential.Id)
	}
	if getErr != nil {
		return nil, errors.New(getErr[0].Summary)
	} else {
		return []*schema.ResourceData{d}, nil
	}
}
