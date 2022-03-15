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

func resourceGoogleCostCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGoogleCostCredentialsCreate,
		ReadContext:   resourceGoogleCostCredentialsRead,
		DeleteContext: resourceGoogleCostCredentialsDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceGcpCredentialsImport},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name for the credentials",
				Required:    true,
			},
			"table_Id": {
				Type:        schema.TypeString,
				Description: "the table id of this credentials ",
				Required:    true,
			},
			"secret": {
				Type:        schema.TypeString,
				Description: "the secret of this credentials",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceGoogleCostCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	value := client.GoogleCostCredentialsValeuPayload{
		TableId: d.Get("table_Id").(string),
		Secret:  d.Get("secret").(string),
	}

	request := client.GoogleCostCredentialsCreatePayload{
		Name:  d.Get("name").(string),
		Value: value,
		Type:  client.GoogleCostCredentiassType,
	}
	credentials, err := apiClient.GoogleCostCredentialsCreate(request)
	if err != nil {
		return diag.Errorf("could not create credentials: %v", err)
	}

	d.SetId(credentials.Id)

	return nil
}

func resourceGoogleCostCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	_, err := apiClient.CloudCredentials(id)
	if err != nil {
		return diag.Errorf("could not get credentials: %v", err)
	}
	return nil
}

func resourceGoogleCostCredentialsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	err := apiClient.CloudCredentialsDelete(id)
	if err != nil {
		return diag.Errorf("could not delete credentials: %v", err)
	}
	return nil
}

func resourceGoogleCostCredentialsImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()
	var getErr diag.Diagnostics
	_, uuidErr := uuid.Parse(id)
	if uuidErr == nil {
		log.Println("[INFO] Resolving GCP Credentials by id: ", id)
		_, getErr = getGoogleCostCredentialsById(id, meta)
	} else {
		log.Println("[DEBUG] ID is not a valid env0 id ", id)
		log.Println("[INFO] Resolving GCP Credentials by name: ", id)
		var gcpCredential client.ApiKey
		gcpCredential, getErr = getGoogleCostCredentialsByName(id, meta)
		d.SetId(gcpCredential.Id)
	}
	if getErr != nil {
		return nil, errors.New(getErr[0].Summary)
	} else {
		return []*schema.ResourceData{d}, nil
	}
}
