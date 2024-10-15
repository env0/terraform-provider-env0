package env0

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSshKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSshKeyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "the name of the ssh key",
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "id of the ssh key",
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
		},
	}
}

func dataSshKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var sshKey *client.SshKey
	var err error

	if name, ok := d.GetOk("name"); ok {
		sshKey, err = getSshKeyByName(name, meta)
		if err != nil {
			return diag.Errorf("could not read ssh key: %v", err)
		}
	} else {
		id := d.Get("id")
		sshKey, err = getSshKeyById(id, meta)
		if err != nil {
			return diag.Errorf("could not read ssh key: %v", err)
		}
		if sshKey == nil {
			return diag.Errorf("could not read ssh key: id %s not found", id)
		}
	}

	if err := writeResourceData(sshKey, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func getSshKeyByName(name interface{}, meta interface{}) (*client.SshKey, error) {
	apiClient := meta.(client.ApiClientInterface)

	return backoff.RetryWithData(func() (*client.SshKey, error) {
		sshKeys, err := apiClient.SshKeys()
		if err != nil {
			return nil, err
		}

		var sshKeysByName []client.SshKey
		for _, candidate := range sshKeys {
			if candidate.Name == name {
				sshKeysByName = append(sshKeysByName, candidate)
			}
		}

		if len(sshKeysByName) > 1 {
			return nil, backoff.Permanent(fmt.Errorf("found multiple ssh keys with name: %s. Use id instead or make sure ssh key names are unique %v", name, sshKeysByName))
		}

		if len(sshKeysByName) == 0 {
			return nil, fmt.Errorf("ssh key with name %v not found", name)
		}

		return &sshKeysByName[0], nil
	}, backoff.NewExponentialBackOff(backoff.WithMaxElapsedTime(time.Minute*1), backoff.WithMaxInterval(time.Second*10)))
}

func getSshKeyById(id interface{}, meta interface{}) (*client.SshKey, error) {
	apiClient := meta.(client.ApiClientInterface)

	sshKeys, err := apiClient.SshKeys()
	if err != nil {
		return nil, err
	}

	var sshKey *client.SshKey

	for _, candidate := range sshKeys {
		if candidate.Id == id.(string) {
			sshKey = &candidate
		}
	}

	return sshKey, nil
}
