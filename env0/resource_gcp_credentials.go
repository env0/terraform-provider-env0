package env0

import (
	"context"
	"errors"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceGcpCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGcpCredentialsCreate,
		ReadContext:   resourceGcpCredentialsRead,
		DeleteContext: resourceGcpCredentialsDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceGcpCredentialsImport},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name for the credentials",
				Required:    true,
				ForceNew:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "the gcp project id",
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
			},
			"service_account_key": {
				Type:        schema.TypeString,
				Description: "the gcp service account key",
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
			},
		},
	}
}

func resourceGcpCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	value := client.GcpCredentialsValuePayload{
		ProjectId:         d.Get("project_id").(string),
		ServiceAccountKey: d.Get("service_account_key").(string),
	}
	requestType := client.GcpServiceAccountCredentialsType

	request := client.GcpCredentialsCreatePayload{
		Name:  d.Get("name").(string),
		Value: value,
		Type:  requestType,
	}
	credentials, err := apiClient.GcpCredentialsCreate(request)
	if err != nil {
		return diag.Errorf("could not create credentials key: %v", err)
	}

	d.SetId(credentials.Id)

	return nil
}

func resourceGcpCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	_, err := apiClient.CloudCredentials(id)
	if err != nil {
		return diag.Errorf("could not get credentials: %v", err)
	}
	return nil
}

func resourceGcpCredentialsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	err := apiClient.CloudCredentialsDelete(id)
	if err != nil {
		return diag.Errorf("could not delete credentials: %v", err)
	}
	return nil
}

func resourceGcpCredentialsImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()
	var getErr diag.Diagnostics
	_, uuidErr := uuid.Parse(id)
	if uuidErr == nil {
		log.Println("[INFO] Resolving GCP Credentials by id: ", id)
		_, getErr = getGcpCredentialsById(id, meta)
	} else {
		log.Println("[DEBUG] ID is not a valid env0 id ", id)
		log.Println("[INFO] Resolving GCP Credentials by name: ", id)
		var gcpCredential client.ApiKey
		gcpCredential, getErr = getGcpCredentialsByName(id, meta)
		d.SetId(gcpCredential.Id)
	}
	if getErr != nil {
		return nil, errors.New(getErr[0].Summary)
	} else {
		return []*schema.ResourceData{d}, nil
	}
}
