package env0

import (
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestReadResourceDataModule(t *testing.T) {
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

func TestReadResourceDataNotification(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceNotification().Schema, map[string]interface{}{
		"name":  "name",
		"type":  "Slack",
		"value": "value",
	})

	expectedPayload := client.NotificationCreatePayload{
		Name:  "name",
		Type:  "Slack",
		Value: "value",
	}

	var payload client.NotificationCreatePayload

	assert.Nil(t, readResourceData(&payload, d))
	assert.Equal(t, expectedPayload, payload)
}

func TestReadResourceDataWithTag(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceAwsCredentials().Schema, map[string]interface{}{
		"name":        "name",
		"arn":         "tagged_arn",
		"external_id": "external_id",
	})

	expectedPayload := client.AwsCredentialsValuePayload{
		RoleArn:    "tagged_arn",
		ExternalId: "external_id",
	}

	var payload client.AwsCredentialsValuePayload

	assert.Nil(t, readResourceData(&payload, d))
	assert.Equal(t, expectedPayload, payload)
}

func TestWriteResourceDataNotification(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceNotification().Schema, map[string]interface{}{})

	n := client.Notification{
		Id:    "id",
		Name:  "name",
		Type:  "Teams",
		Value: "value",
	}

	assert.Nil(t, writeResourceData(&n, d))

	assert.Equal(t, "id", d.Id())
	assert.Equal(t, "name", d.Get("name"))
	assert.Equal(t, "Teams", d.Get("type"))
	assert.Equal(t, "value", d.Get("value"))
}

func TestReadResourceDataNotificationProjectAssignment(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceNotificationProjectAssignment().Schema, map[string]interface{}{
		"event_names": []interface{}{
			"driftUndetected",
			"destroySucceeded",
		},
	})

	expectedPayload := client.NotificationProjectAssignmentUpdatePayload{
		EventNames: []string{
			"driftUndetected",
			"destroySucceeded",
		},
	}

	var payload client.NotificationProjectAssignmentUpdatePayload
	assert.Nil(t, readResourceData(&payload, d))
	assert.Equal(t, expectedPayload, payload)
}

func TestWriteResourceDataNotificationProjectAssignment(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceNotificationProjectAssignment().Schema, map[string]interface{}{})

	a := client.NotificationProjectAssignment{
		Id:                     "id",
		NotificationEndpointId: "nid",
		EventNames: []string{
			"driftUndetected",
		},
	}

	assert.Nil(t, writeResourceData(&a, d))

	assert.Equal(t, "id", d.Id())
	assert.Equal(t, "nid", d.Get("notification_endpoint_id"))
	assert.Equal(t, []interface{}{"driftUndetected"}, d.Get("event_names"))
}

func TestWriteCustomResourceData(t *testing.T) {
	d := schema.TestResourceDataRaw(t, dataConfigurationVariable().Schema, map[string]interface{}{})

	configurationVariable := client.ConfigurationVariable{
		Id:             "id0",
		Name:           "name0",
		Description:    "desc0",
		ScopeId:        "scope0",
		Value:          "value0",
		OrganizationId: "organization0",
		UserId:         "user0",
		IsSensitive:    boolPtr(true),
		Scope:          client.ScopeEnvironment,
		Type:           (*client.ConfigurationVariableType)(intPtr(1)),
		Schema:         &client.ConfigurationVariableSchema{Type: "string", Format: client.HCL, Enum: []string{"s1", "s2"}},
		IsReadOnly:     boolPtr(true),
		IsRequired:     boolPtr(true),
		ToDelete:       boolPtr(false),
		Regex:          "regex",
	}

	assert.Nil(t, writeResourceData(&configurationVariable, d))

	assert.Equal(t, configurationVariable.Id, d.Id())
	assert.Equal(t, configurationVariable.Name, d.Get("name"))
	assert.Equal(t, configurationVariable.Description, d.Get("description"))
	assert.Equal(t, "terraform", d.Get("type"))
	assert.Equal(t, configurationVariable.Value, d.Get("value"))
	assert.Equal(t, string(configurationVariable.Scope), d.Get("scope"))
	assert.Equal(t, *configurationVariable.IsReadOnly, d.Get("is_read_only"))
	assert.Equal(t, *configurationVariable.IsRequired, d.Get("is_required"))
	assert.Equal(t, configurationVariable.Regex, d.Get("regex"))
}

func TestReadCustomResourceData(t *testing.T) {
	d := schema.TestResourceDataRaw(t, dataConfigurationVariable().Schema, map[string]interface{}{
		"type":        "terraform",
		"name":        "name",
		"description": "description",
	})

	params := client.ConfigurationVariableCreateParams{}

	assert.Nil(t, readResourceData(&params, d))

	assert.Equal(t, params.Name, "name")
	assert.Equal(t, int(params.Type), 1)
	assert.Equal(t, params.Description, "description")
}

func TestWriteResourceDataSliceVariablesAgents(t *testing.T) {
	d := schema.TestResourceDataRaw(t, dataAgents().Schema, map[string]interface{}{})

	agent1 := client.Agent{
		AgentKey: "key1",
	}

	agent2 := client.Agent{
		AgentKey: "key1",
	}

	vars := []client.Agent{agent1, agent2}

	assert.Nil(t, writeResourceDataSlice(vars, "agents", d))
	assert.Equal(t, agent1.AgentKey, d.Get("agents.0.agent_key"))
	assert.Equal(t, agent2.AgentKey, d.Get("agents.1.agent_key"))
}

func TestWriteResourceDataSliceVariablesConfigurationVariable(t *testing.T) {
	d := schema.TestResourceDataRaw(t, dataSourceCodeVariables().Schema, map[string]interface{}{})

	schema1 := client.ConfigurationVariableSchema{
		Type:   "string",
		Format: "HCL",
		Enum:   []string{"Variable", "a"},
	}

	schema2 := client.ConfigurationVariableSchema{
		Type:   "string",
		Format: "JSON",
	}

	var1 := client.ConfigurationVariable{
		Id:          "id0",
		Name:        "name0",
		Description: "desc0",
		Value:       "v1",
		Schema:      &schema1,
	}

	var2 := client.ConfigurationVariable{
		Id:          "id1",
		Name:        "name1",
		Description: "desc1",
		Value:       "v2",
		Schema:      &schema2,
	}

	vars := []client.ConfigurationVariable{var1, var2}

	assert.Nil(t, writeResourceDataSlice(vars, "variables", d))
	assert.Equal(t, var1.Name, d.Get("variables.0.name"))
	assert.Equal(t, var2.Name, d.Get("variables.1.name"))
	assert.Equal(t, var1.Value, d.Get("variables.0.value"))
	assert.Equal(t, var2.Value, d.Get("variables.1.value"))
	assert.Equal(t, string(var1.Schema.Format), d.Get("variables.0.format"))
	assert.Equal(t, string(var2.Schema.Format), d.Get("variables.1.format"))
}
