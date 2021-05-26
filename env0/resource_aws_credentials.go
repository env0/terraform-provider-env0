package env0

import (
	"context"
	"encoding/json"
	"errors"

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

func resourceSshKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*client.ApiClient)

	sshKeys, err := apiClient.SshKeys()
	if err != nil {
		return diag.Errorf("could not query ssh keys: %v", err)
	}
	found := false
	for _, candidate := range sshKeys {
		if candidate.Id == d.Id() {
			found = true
		}
	}
	if !found {
		return diag.Errorf("ssh key %s not found", d.Id())
	}

	return nil
}

func resourceSshKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*client.ApiClient)

	id := d.Id()
	err := apiClient.SshKeyDelete(id)
	if err != nil {
		return diag.Errorf("could not delete ssh key: %v", err)
	}
	return nil
}

func resourceSshKeyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return nil, errors.New("Not implemented")
	// apiClient := meta.(*client.ApiClient)

	// id := d.Id()
	// ssh key, err := apiClient.SshKey(id)
	// if err != nil {
	// 	return nil, err
	// }

	// d.Set("name", ssh key.Name)

	// return []*schema.ResourceData{d}, nil
}
