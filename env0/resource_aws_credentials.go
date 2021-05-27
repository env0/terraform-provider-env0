package env0

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"

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
			},
			"arn": {
				Type:        schema.TypeString,
				Description: "the aws role arn",
				Required:    true,
			},
			"external_id": {
				Type:        schema.TypeString,
				Description: "the aws role external id",
				Required:    true,
			},
		},
	}
}

func resourceAwsCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*client.ApiClient)
	value, _ := json.Marshal(map[string]string{
		"arn":        d.Get("arn").(string),
		"externalId": d.Get("external_id").(string),
	})
	request := client.AwsCredentialsCreatePayload{
		Name:  d.Get("name").(string),
		Type:  "aws_assumed_role",
		Value: value,
	}
	credentials, err := apiClient.AwsCredentialsCreate(request)
	if err != nil {
		return diag.Errorf("could not create credentials key: %v", err)
	}

	d.SetId(credentials.Id)

	return nil
}

func resourceAwsCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*client.ApiClient)

	id := d.Id()
	_, err := apiClient.AwsCredentials(id)
	if err != nil {
		return diag.Errorf("could not get credentials: %v", err)
	}
	return nil
}

func resourceAwsCredentialsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*client.ApiClient)

	id := d.Id()
	err := apiClient.AwsCredentialsDelete(id)
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
		var project client.Project
		project, getErr = getAwsCredentialsByName(id, meta)
		d.SetId(project.Id)
	}
	if getErr != nil {
		return nil, errors.New(getErr[0].Summary)
	} else {
		return []*schema.ResourceData{d}, nil
	}
}
