package env0

import (
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestReadResourceDataModule(t *testing.T) {
	t.Run("match", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceModule().Schema, map[string]interface{}{
			"module_name":            "module_name",
			"module_provider":        "module_provider",
			"github_installation_id": 1000,
			"ssh_keys": []interface{}{
				map[string]interface{}{"id": "id1", "name": "name1"},
			},
		})

		githubInstallationId := 1000

		expectedPayload := client.ModuleCreatePayload{
			ModuleName:           "module_name",
			ModuleProvider:       "module_provider",
			GithubInstallationId: &githubInstallationId,
			SshKeys: []client.ModuleSshKey{
				{
					Id:   "id1",
					Name: "name1",
				},
			},
		}

		var payload client.ModuleCreatePayload

		assert.Nil(t, readResourceData(&payload, d))
		assert.Equal(t, expectedPayload, payload)
	})
}

func TestWriteResourceDataModule(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceModule().Schema, map[string]interface{}{})

	m := client.Module{
		ModuleName:           "module_name",
		BitbucketClientKey:   stringPtr("12345abcd"),
		ModuleProvider:       "module_provider",
		Repository:           "repository",
		Description:          "description",
		Id:                   "id",
		GithubInstallationId: intPtr(1000),
		SshKeys: []client.ModuleSshKey{
			{
				Id:   "id1",
				Name: "name1",
			},
		},
	}

	assert.Nil(t, writeResourceData(&m, d))

	assert.Equal(t, "id", d.Id())
	assert.Equal(t, "module_name", d.Get("module_name"))
	assert.Equal(t, "12345abcd", d.Get("bitbucket_client_key"))
	assert.Equal(t, "module_provider", d.Get("module_provider"))
	assert.Equal(t, "repository", d.Get("repository"))
	assert.Equal(t, "description", d.Get("description"))
	assert.Equal(t, 1000, d.Get("github_installation_id"))

	var rawSshKeys []interface{}
	rawSshKeys = append(rawSshKeys, map[string]interface{}{"id": "id1", "name": "name1"})
	assert.Equal(t, rawSshKeys, d.Get("ssh_keys"))
}
