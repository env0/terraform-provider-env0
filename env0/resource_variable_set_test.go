package env0

import (
	"fmt"
	"regexp"
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
	projectId := "proj"

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

	senstiveTextVariable := client.ConfigurationVariable{
		Name:           "nv1",
		Value:          "v2",
		OrganizationId: organizationId,
		IsSensitive:    boolPtr(true),
		Scope:          "SET",
		Type:           (*client.ConfigurationVariableType)(intPtr(0)),
		Schema: &client.ConfigurationVariableSchema{
			Type: "string",
		},
	}

	senstiveTextVariableWithScopeId := senstiveTextVariable
	senstiveTextVariableWithScopeId.Id = "idsentivetextvariable"
	senstiveTextVariableWithScopeId.ScopeId = configurationSet.Id
	senstiveTextVariableWithScopeId.Value = "OMITTED"

	hclVariable := client.ConfigurationVariable{
		Name:           "hcl1",
		Value:          "sdzdfsdfsd",
		OrganizationId: organizationId,
		IsSensitive:    boolPtr(false),
		Scope:          "SET",
		Type:           (*client.ConfigurationVariableType)(intPtr(1)),
		Schema: &client.ConfigurationVariableSchema{
			Format: "HCL",
		},
	}

	hclVariableWithScopeId := hclVariable
	hclVariableWithScopeId.Id = "idhclvariable"
	hclVariableWithScopeId.ScopeId = configurationSet.Id

	jsonVariable := client.ConfigurationVariable{
		Name:           "json1",
		Value:          "{}",
		OrganizationId: organizationId,
		IsSensitive:    boolPtr(false),
		Scope:          "SET",
		Type:           (*client.ConfigurationVariableType)(intPtr(1)),
		Schema: &client.ConfigurationVariableSchema{
			Format: "JSON",
		},
	}

	jsonVariableWithScopeId := jsonVariable
	jsonVariableWithScopeId.Id = "idjsonvariable"
	jsonVariableWithScopeId.ScopeId = configurationSet.Id

	dropdownVariable := client.ConfigurationVariable{
		Name:           "dropdown123",
		Value:          "o1",
		OrganizationId: organizationId,
		IsSensitive:    boolPtr(false),
		Scope:          "SET",
		Type:           (*client.ConfigurationVariableType)(intPtr(1)),
		Schema: &client.ConfigurationVariableSchema{
			Type: "string",
			Enum: []string{
				"o1", "o2",
			},
		},
	}

	dropdownVariableWithScopeId := dropdownVariable
	dropdownVariableWithScopeId.Id = "iddropdownvariable"
	dropdownVariableWithScopeId.ScopeId = configurationSet.Id

	t.Run("basic - organization scope", func(t *testing.T) {
		createPayload := client.CreateConfigurationSetPayload{
			Name:        configurationSet.Name,
			Description: configurationSet.Description,
			Scope:       "organization",
			ConfigurationProperties: []client.ConfigurationVariable{
				textVariable,
				senstiveTextVariable,
				hclVariable,
				jsonVariable,
				dropdownVariable,
			},
		}

		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
					resource "%s" "%s" {
						name = "%s"
						description = "%s"

						variable {
							name = "%s"
							value = "%s"
							type = "terraform"
							format = "text"
						}

						variable {
							name = "%s"
							value = "%s"
							type = "environment"
							format = "text"
							is_sensitive = true
						}

						variable {
							name = "%s"
							value = "%s"
							type = "terraform"
							format = "hcl"
						}

						variable {
							name = "%s"
							value = "%s"
							type = "terraform"
							format = "json"
						}

						variable {
							name = "%s"
							dropdown_values = ["o1", "o2"]
							type = "terraform"
							format = "dropdown"
						}
					}`, resourceType, resourceName, createPayload.Name, createPayload.Description,
						textVariable.Name, textVariable.Value,
						senstiveTextVariable.Name, senstiveTextVariable.Value,
						hclVariable.Name, hclVariable.Value,
						jsonVariable.Name, jsonVariable.Value,
						dropdownVariable.Name,
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", configurationSet.Id),
						resource.TestCheckResourceAttr(accessor, "name", configurationSet.Name),
						resource.TestCheckResourceAttr(accessor, "description", configurationSet.Description),
						resource.TestCheckResourceAttr(accessor, "scope", "organization"),
						resource.TestCheckResourceAttr(accessor, "variable.0.value", textVariable.Value),
						resource.TestCheckResourceAttr(accessor, "variable.0.name", textVariable.Name),
						resource.TestCheckResourceAttr(accessor, "variable.0.type", "terraform"),
						resource.TestCheckResourceAttr(accessor, "variable.0.format", "text"),
						resource.TestCheckResourceAttr(accessor, "variable.1.value", senstiveTextVariable.Value),
						resource.TestCheckResourceAttr(accessor, "variable.1.name", senstiveTextVariable.Name),
						resource.TestCheckResourceAttr(accessor, "variable.1.type", "environment"),
						resource.TestCheckResourceAttr(accessor, "variable.1.format", "text"),
						resource.TestCheckResourceAttr(accessor, "variable.1.is_sensitive", "true"),
						resource.TestCheckResourceAttr(accessor, "variable.2.value", hclVariable.Value),
						resource.TestCheckResourceAttr(accessor, "variable.2.name", hclVariable.Name),
						resource.TestCheckResourceAttr(accessor, "variable.2.type", "terraform"),
						resource.TestCheckResourceAttr(accessor, "variable.2.format", "hcl"),
						resource.TestCheckResourceAttr(accessor, "variable.3.value", jsonVariable.Value),
						resource.TestCheckResourceAttr(accessor, "variable.3.name", jsonVariable.Name),
						resource.TestCheckResourceAttr(accessor, "variable.3.type", "terraform"),
						resource.TestCheckResourceAttr(accessor, "variable.3.format", "json"),
						resource.TestCheckResourceAttr(accessor, "variable.4.dropdown_values.0", dropdownVariable.Schema.Enum[0]),
						resource.TestCheckResourceAttr(accessor, "variable.4.dropdown_values.1", dropdownVariable.Schema.Enum[1]),
						resource.TestCheckResourceAttr(accessor, "variable.4.name", dropdownVariable.Name),
						resource.TestCheckResourceAttr(accessor, "variable.4.type", "terraform"),
						resource.TestCheckResourceAttr(accessor, "variable.4.format", "dropdown"),
					),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().OrganizationId().AnyTimes().Return(organizationId, nil)

			gomock.InOrder(
				mock.EXPECT().ConfigurationSetCreate(&createPayload).Times(1).Return(&configurationSet, nil),
				mock.EXPECT().ConfigurationSet(configurationSet.Id).Times(1).Return(&configurationSet, nil),
				mock.EXPECT().ConfigurationVariablesBySetId(configurationSet.Id).Times(1).Return([]client.ConfigurationVariable{
					textVariableWithScopeId,
					senstiveTextVariableWithScopeId,
					hclVariableWithScopeId,
					jsonVariableWithScopeId,
					dropdownVariableWithScopeId,
				}, nil),
				mock.EXPECT().ConfigurationSetDelete(configurationSet.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("basic - project scope", func(t *testing.T) {
		createPayload := client.CreateConfigurationSetPayload{
			Name:        configurationSet.Name,
			Description: configurationSet.Description,
			Scope:       "project",
			ScopeId:     projectId,
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
						scope = "project"
						scope_id = "%s"

						variable {
							name = "%s"
							value = "%s"
							type = "terraform"
							format = "text"
						}
					}`, resourceType, resourceName, createPayload.Name, createPayload.Description, createPayload.ScopeId,
						textVariable.Name, textVariable.Value,
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", configurationSet.Id),
						resource.TestCheckResourceAttr(accessor, "name", configurationSet.Name),
						resource.TestCheckResourceAttr(accessor, "description", configurationSet.Description),
						resource.TestCheckResourceAttr(accessor, "scope", "project"),
						resource.TestCheckResourceAttr(accessor, "scope_id", projectId),
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
				mock.EXPECT().ConfigurationVariablesBySetId(configurationSet.Id).Times(1).Return([]client.ConfigurationVariable{
					textVariableWithScopeId,
				}, nil),
				mock.EXPECT().ConfigurationSetDelete(configurationSet.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("update", func(t *testing.T) {
		createPayload := client.CreateConfigurationSetPayload{
			Name:        configurationSet.Name,
			Description: configurationSet.Description,
			Scope:       "organization",
			ConfigurationProperties: []client.ConfigurationVariable{
				textVariable,
			},
		}

		updatedTextVariable := textVariable
		updatedTextVariable.Value = "new-value"

		updatedTextVariableWithScopeId := textVariableWithScopeId
		updatedTextVariableWithScopeId.Value = updatedTextVariable.Value

		updateTestCase := resource.TestCase{
			// Create a text variable.
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
					resource "%s" "%s" {
						name = "%s"
						description = "%s"

						variable {
							name = "%s"
							value = "%s"
							type = "terraform"
							format = "text"
						}
					}`, resourceType, resourceName, configurationSet.Name, configurationSet.Description,
						textVariable.Name, textVariable.Value,
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", configurationSet.Id),
						resource.TestCheckResourceAttr(accessor, "name", configurationSet.Name),
						resource.TestCheckResourceAttr(accessor, "description", configurationSet.Description),
						resource.TestCheckResourceAttr(accessor, "scope", "organization"),
						resource.TestCheckResourceAttr(accessor, "variable.0.value", textVariable.Value),
						resource.TestCheckResourceAttr(accessor, "variable.0.name", textVariable.Name),
						resource.TestCheckResourceAttr(accessor, "variable.0.type", "terraform"),
						resource.TestCheckResourceAttr(accessor, "variable.0.format", "text"),
					),
				},
				// Update the value of a text variable.
				{
					Config: fmt.Sprintf(`
					resource "%s" "%s" {
						name = "%s"
						description = "%s"

						variable {
							name = "%s"
							value = "%s"
							type = "terraform"
							format = "text"
						}
					}`, resourceType, resourceName, configurationSet.Name, configurationSet.Description,
						updatedTextVariable.Name, updatedTextVariable.Value,
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", configurationSet.Id),
						resource.TestCheckResourceAttr(accessor, "name", configurationSet.Name),
						resource.TestCheckResourceAttr(accessor, "description", configurationSet.Description),
						resource.TestCheckResourceAttr(accessor, "scope", "organization"),
						resource.TestCheckResourceAttr(accessor, "variable.0.value", updatedTextVariable.Value),
						resource.TestCheckResourceAttr(accessor, "variable.0.name", textVariable.Name),
						resource.TestCheckResourceAttr(accessor, "variable.0.type", "terraform"),
						resource.TestCheckResourceAttr(accessor, "variable.0.format", "text"),
					),
				},
			},
		}

		runUnitTest(t, updateTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().OrganizationId().AnyTimes().Return(organizationId, nil)
			mock.EXPECT().ConfigurationSet(configurationSet.Id).AnyTimes().Return(&configurationSet, nil)

			gomock.InOrder(
				mock.EXPECT().ConfigurationSetCreate(&createPayload).Times(1).Return(&configurationSet, nil),
				mock.EXPECT().ConfigurationVariablesBySetId(configurationSet.Id).Times(3).Return([]client.ConfigurationVariable{
					textVariableWithScopeId,
				}, nil),
				mock.EXPECT().ConfigurationSetUpdate(configurationSet.Id, &client.UpdateConfigurationSetPayload{
					Name:                           configurationSet.Name,
					Description:                    configurationSet.Description,
					ConfigurationPropertiesChanges: []client.ConfigurationVariable{updatedTextVariableWithScopeId},
				}),
				mock.EXPECT().ConfigurationVariablesBySetId(configurationSet.Id).Times(1).Return([]client.ConfigurationVariable{
					updatedTextVariableWithScopeId,
				}, nil),
				mock.EXPECT().ConfigurationSetDelete(configurationSet.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("failures", func(t *testing.T) {
		runFailure := func(testName string, config string, errorMessage string) {
			t.Run(testName, func(t *testing.T) {
				testCase := resource.TestCase{
					Steps: []resource.TestStep{
						{
							Config:      config,
							ExpectError: regexp.MustCompile(errorMessage),
						},
					},
				}

				runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
					mock.EXPECT().OrganizationId().AnyTimes().Return(organizationId, nil)
				})
			})
		}

		runFailure("no value - text", fmt.Sprintf(`
			resource "%s" "%s" {
				name = "fail"
				description = "description1"
				scope = "organization"

				variable {
					name = "a"
					type = "terraform"
					format = "text"
				}
			}`, resourceType, resourceName), "free text variable 'a' must have a value")

		runFailure("no value - hcl", fmt.Sprintf(`
			resource "%s" "%s" {
				name = "fail"
				description = "description1"
				scope = "organization"

				variable {
					name = "a"
					type = "terraform"
					format = "hcl"
				}
			}`, resourceType, resourceName), "hcl variable 'a' must have a value")

		runFailure("no value - json", fmt.Sprintf(`
			resource "%s" "%s" {
				name = "fail"
				description = "description1"
				scope = "organization"

				variable {
					name = "a"
					type = "terraform"
					format = "json"
				}
			}`, resourceType, resourceName), "json variable 'a' must have a value")

		runFailure("no dropdown_values", fmt.Sprintf(`
			resource "%s" "%s" {
				name = "fail"
				description = "description1"
				scope = "organization"

				variable {
					name = "a"
					type = "terraform"
					format = "dropdown"
				}
			}`, resourceType, resourceName), "dropdown variable 'a' must have dropdown_values")

		runFailure("invalid json", fmt.Sprintf(`
			resource "%s" "%s" {
				name = "fail"
				description = "description1"
				scope = "organization"

				variable {
					name = "a"
					type = "terraform"
					format = "json"
					value = "i am not a valid json"
				}
			}`, resourceType, resourceName), "json variable 'a' is not a valid json value")

		runFailure("project scope with no scope id", fmt.Sprintf(`
			resource "%s" "%s" {
				name = "fail"
				description = "description1"
				scope = "project"
			}`, resourceType, resourceName), "scope_id must be configured for the scope 'project'")
	})
}
