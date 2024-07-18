package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestVariableSetDataSource(t *testing.T) {
	projectId := "project_id"
	organizationId := "organization_id"

	v1 := client.ConfigurationSet{
		Id:              "id1",
		Name:            "name1",
		CreationScopeId: projectId,
	}

	v2 := client.ConfigurationSet{
		Id:              "id2",
		Name:            "name2",
		CreationScopeId: projectId,
	}

	v3 := client.ConfigurationSet{
		Id:              "id3",
		Name:            "name3",
		CreationScopeId: organizationId,
	}

	v4 := client.ConfigurationSet{
		Id:              "id4",
		Name:            "name4",
		CreationScopeId: "some_other_id",
	}

	resourceType := "env0_variable_set"
	resourceName := "test_variable_set"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getConfig := func(name string, scope string, projectId string) string {
		fields := map[string]interface{}{"name": name, "scope": scope}
		if projectId != "" {
			fields["project_id"] = projectId
		}
		return dataSourceConfigCreate(resourceType, resourceName, fields)
	}

	mockVariableSetsCall := func(organizationId string, projectId string, returnValue []client.ConfigurationSet) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			if organizationId != "" {
				mock.EXPECT().OrganizationId().AnyTimes().Return(organizationId, nil)
			}
			mock.EXPECT().ConfigurationSets(organizationId, projectId).AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("project id scope", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: getConfig(v2.Name, "PROJECT", v2.CreationScopeId),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "id", v2.Id),
						),
					},
				},
			},
			mockVariableSetsCall("", projectId, []client.ConfigurationSet{
				v4, v1, v2,
			}),
		)
	})

	t.Run("organization id scope", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: getConfig(v3.Name, "ORGANIZATION", ""),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "id", v3.Id),
						),
					},
				},
			},
			mockVariableSetsCall(organizationId, "", []client.ConfigurationSet{
				v4, v3,
			}),
		)
	})

	t.Run("name not found", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config:      getConfig("name that isn't found", "PROJECT", v2.CreationScopeId),
						ExpectError: regexp.MustCompile("variable set not found"),
					},
				},
			},
			mockVariableSetsCall("", projectId, []client.ConfigurationSet{
				v4, v1, v2, v3,
			}),
		)
	})

	t.Run("get configuration sets api call failed", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config:      getConfig(v2.Name, "PROJECT", v2.CreationScopeId),
						ExpectError: regexp.MustCompile("could not get variable sets: error"),
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().ConfigurationSets("", projectId).AnyTimes().Return(nil, errors.New("error"))
			},
		)
	})

	t.Run("get organization id api call failed", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config:      getConfig(v3.Name, "ORGANIZATION", ""),
						ExpectError: regexp.MustCompile("could not get organization id: error"),
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().OrganizationId().AnyTimes().Return("", errors.New("error"))
			},
		)
	})

	t.Run("project_id is required", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config:      getConfig(v2.Name, "PROJECT", ""),
						ExpectError: regexp.MustCompile("'project_id' is required"),
					},
				},
			},
			func(mock *client.MockApiClientInterface) {},
		)
	})
}
