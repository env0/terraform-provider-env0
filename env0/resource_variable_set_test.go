package env0

import (
	"fmt"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitVariableSetResource(t *testing.T) {
	resourceType := "env0_variable_set"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	organizationId := "org"

	configurationSet := client.ConfigurationSet{
		Id:          "idddd111",
		Name:        "name1",
		Description: "des1",
	}

	textVariable := client.ConfigurationVariable{
		Name:           "nv1",
		Value:          "v1",
		OrganizationId: organizationId,
		IsSensitive:    boolPtr(false),
		Scope:          "SET",
		Type:           (*client.ConfigurationVariableType)(intPtr(1)),
		Schema: &client.ConfigurationVariableSchema{
			Type: "string",
		},
	}

	textVariableWithScopeId := textVariable
	textVariableWithScopeId.Id = "idtextvariable"
	textVariableWithScopeId.ScopeId = configurationSet.Id

	t.Run("basic - organization scope", func(t *testing.T) {
		createPayload := client.CreateConfigurationSetPayload{
			Name:        configurationSet.Name,
			Description: configurationSet.Description,
			Scope:       "organization",
			ConfigurationProperties: []client.ConfigurationVariable{
				textVariable,
			},
		}

		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
					resource "%s" "%s" {
						name = "%s"
						description = "%s"
						scope = "organization"

						variable {
							name = "%s"
							value = "%s"
							type = "terraform"
							format = "text"
						}
					}`, resourceType, resourceName, createPayload.Name, createPayload.Description,
						textVariable.Name, textVariable.Value,
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", configurationSet.Id),
						resource.TestCheckResourceAttr(accessor, "name", configurationSet.Name),
						resource.TestCheckResourceAttr(accessor, "description", configurationSet.Description),
						resource.TestCheckResourceAttr(accessor, "variable.0.value", textVariable.Value),
						resource.TestCheckResourceAttr(accessor, "variable.0.name", textVariable.Name),
						resource.TestCheckResourceAttr(accessor, "variable.0.type", "terraform"),
						resource.TestCheckResourceAttr(accessor, "variable.0.format", "text"),
					),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().OrganizationId().AnyTimes().Return(organizationId, nil)

			gomock.InOrder(
				mock.EXPECT().ConfigurationSetCreate(&createPayload).Times(1).Return(&configurationSet, nil),
				mock.EXPECT().ConfigurationSet(configurationSet.Id).Times(1).Return(&configurationSet, nil),
				mock.EXPECT().ConfigurationVariablesBySetId(configurationSet.Id).Times(1).Return([]client.ConfigurationVariable{textVariableWithScopeId}, nil),
				mock.EXPECT().ConfigurationSetDelete(configurationSet.Id).Times(1).Return(nil),
			)
		})
	})

}
