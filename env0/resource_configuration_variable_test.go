package env0

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"text/template"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitConfigurationVariableResource(t *testing.T) {
	resourceType := "env0_configuration_variable"
	resourceName := "test"
	isReadonly := true
	isRequired := false
	accessor := resourceAccessor(resourceType, resourceName)
	configVar := client.ConfigurationVariable{
		Id:          "id0",
		Name:        "name0",
		Description: "desc0",
		Value:       "Variable",
		IsReadOnly:  &isReadonly,
		IsRequired:  &isRequired,
	}
	stepConfig := resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
		"name":         configVar.Name,
		"description":  configVar.Description,
		"value":        configVar.Value,
		"is_read_only": strconv.FormatBool(*configVar.IsReadOnly),
		"is_required":  strconv.FormatBool(*configVar.IsRequired),
	})

	configurationVariableCreateParams := client.ConfigurationVariableCreateParams{
		Name:        configVar.Name,
		Value:       configVar.Value,
		IsSensitive: false,
		Scope:       client.ScopeGlobal,
		ScopeId:     "",
		Type:        client.ConfigurationVariableTypeEnvironment,
		EnumValues:  nil,
		Description: configVar.Description,
		Format:      client.Text,
		IsRequired:  *configVar.IsRequired,
		IsReadOnly:  *configVar.IsReadOnly,
	}
	t.Run("Create", func(t *testing.T) {
		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", configVar.Id),
						resource.TestCheckResourceAttr(accessor, "name", configVar.Name),
						resource.TestCheckResourceAttr(accessor, "description", configVar.Description),
						resource.TestCheckResourceAttr(accessor, "value", configVar.Value),
						resource.TestCheckResourceAttr(accessor, "is_read_only", strconv.FormatBool(*configVar.IsReadOnly)),
						resource.TestCheckResourceAttr(accessor, "is_required", strconv.FormatBool(*configVar.IsRequired)),
					),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ConfigurationVariableCreate(configurationVariableCreateParams).Times(1).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariablesById(configVar.Id).Times(1).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)
		})
	})

	t.Run("Create Two with readonly", func(t *testing.T) {
		// https://github.com/env0/terraform-provider-env0/issues/215
		// we want to create two variables, one org level with read only and another in lower level and see we can still manage both - double apply and destroy
		orgResourceName := "org"
		projResourceName := "project"
		orgAccessor := resourceAccessor(resourceType, orgResourceName)
		projAccessor := resourceAccessor(resourceType, projResourceName)

		orgVar := client.ConfigurationVariable{
			Id:         "orgVariableId",
			Name:       "variable",
			Value:      "orgVariable",
			IsReadOnly: &isReadonly,
		}

		orgConfigVariableCreateParams := client.ConfigurationVariableCreateParams{
			Name:        orgVar.Name,
			Value:       orgVar.Value,
			IsSensitive: false,
			Scope:       client.ScopeGlobal,
			ScopeId:     "",
			Type:        client.ConfigurationVariableTypeEnvironment,
			EnumValues:  nil,
			Format:      client.Text,
			IsReadOnly:  *orgVar.IsReadOnly,
		}

		projVar := client.ConfigurationVariable{
			Id:      "projectVariableId",
			Name:    orgVar.Name,
			Value:   "projVariable",
			Scope:   client.ScopeProject,
			ScopeId: "projectId",
		}

		projectConfigVariableCreateParams := client.ConfigurationVariableCreateParams{
			Name:        projVar.Name,
			Value:       projVar.Value,
			IsSensitive: false,
			Scope:       client.ScopeProject,
			ScopeId:     projVar.ScopeId,
			Type:        client.ConfigurationVariableTypeEnvironment,
			EnumValues:  nil,
			Format:      client.Text,
		}

		data := map[string]interface{}{"projectId": projVar.ScopeId, "orgResourceName": orgResourceName, "projResourceName": projResourceName, "resourceType": resourceType, "variableName": orgVar.Name, "orgValue": orgVar.Value, "projValue": projVar.Value}

		tmpl, err := template.New("").Parse(`
resource "{{.resourceType}}" "{{.orgResourceName}}" {
  name = "{{.variableName}}"
  value = "{{.orgValue}}"
  is_read_only = true
}

resource "{{.resourceType}}" "{{.projResourceName}}" {
  name = {{.resourceType}}.{{.orgResourceName}}.name
  value = "{{.projValue}}"
  project_id = "{{.projectId}}"
}
`)
		if err != nil {
			panic(err)
		}
		var tpl bytes.Buffer

		err = tmpl.Execute(&tpl, data)
		if err != nil {
			panic(err)
		}
		stepConfig := tpl.String()

		testStep := resource.TestStep{
			Config: stepConfig,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(orgAccessor, "id", orgVar.Id),
				resource.TestCheckResourceAttr(orgAccessor, "name", orgVar.Name),
				resource.TestCheckResourceAttr(orgAccessor, "value", orgVar.Value),
				resource.TestCheckResourceAttr(orgAccessor, "is_read_only", strconv.FormatBool(*orgVar.IsReadOnly)),
				resource.TestCheckResourceAttr(projAccessor, "id", projVar.Id),
				resource.TestCheckResourceAttr(projAccessor, "name", projVar.Name),
				resource.TestCheckResourceAttr(projAccessor, "value", projVar.Value),
			),
		}
		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				testStep,
				testStep,
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ConfigurationVariableCreate(projectConfigVariableCreateParams).Times(1).Return(projVar, nil)
			mock.EXPECT().ConfigurationVariableCreate(orgConfigVariableCreateParams).Times(1).Return(orgVar, nil)
			mock.EXPECT().ConfigurationVariablesById(orgVar.Id).AnyTimes().Return(orgVar, nil)
			mock.EXPECT().ConfigurationVariablesById(projVar.Id).AnyTimes().Return(projVar, nil)
			mock.EXPECT().ConfigurationVariableDelete(orgVar.Id).Times(1).Return(nil)
			mock.EXPECT().ConfigurationVariableDelete(projVar.Id).Times(1).Return(nil)
		})
	})

	t.Run("Create and update with regex", func(t *testing.T) {
		initialVar := client.ConfigurationVariable{
			Id:    "regex-var-id",
			Name:  "regex-var-name",
			Regex: "initial-regex",
		}
		initialResource := resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
			"name":  initialVar.Name,
			"regex": initialVar.Regex,
		})
		createParams := client.ConfigurationVariableCreateParams{
			Name:  initialVar.Name,
			Regex: initialVar.Regex,
			Scope: client.ScopeGlobal,
		}

		updatedVar := initialVar
		updatedVar.Regex = "updated-regex"
		updatedResource := resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
			"name":  updatedVar.Name,
			"regex": updatedVar.Regex,
		})
		updateParams := createParams
		updateParams.Regex = updatedVar.Regex

		steps := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: initialResource,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "name", initialVar.Name),
						resource.TestCheckResourceAttr(accessor, "regex", initialVar.Regex),
					),
				},
				{
					Config: updatedResource,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "name", updatedVar.Name),
						resource.TestCheckResourceAttr(accessor, "regex", updatedVar.Regex),
					),
				},
			},
		}

		runUnitTest(t, steps, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ConfigurationVariableCreate(createParams).Times(1).Return(initialVar, nil),
				mock.EXPECT().ConfigurationVariablesById(initialVar.Id).Times(2).Return(initialVar, nil),
				mock.EXPECT().ConfigurationVariableUpdate(client.ConfigurationVariableUpdateParams{CommonParams: updateParams, Id: updatedVar.Id}).Times(1).Return(updatedVar, nil),
				mock.EXPECT().ConfigurationVariablesById(initialVar.Id).Times(1).Return(updatedVar, nil),
				mock.EXPECT().ConfigurationVariableDelete(initialVar.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Create Enum", func(t *testing.T) {
		schema := client.ConfigurationVariableSchema{
			Type: "string",
			Enum: []string{"Variable", "a"},
		}
		configVar := client.ConfigurationVariable{
			Id:          "id0",
			Name:        "name0",
			Description: "desc0",
			Value:       "Variable",
			Schema:      &schema,
		}
		stepConfig := fmt.Sprintf(`
	resource "%s" "test" {
		name = "%s"
		description = "%s"
		value= "%s"
		enum = ["%s","%s"]
	}`, resourceType, configVar.Name, configVar.Description, configVar.Value, configVar.Schema.Enum[0], configVar.Schema.Enum[1])

		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", configVar.Id),
						resource.TestCheckResourceAttr(accessor, "name", configVar.Name),
						resource.TestCheckResourceAttr(accessor, "description", configVar.Description),
						resource.TestCheckResourceAttr(accessor, "value", configVar.Value),
						resource.TestCheckResourceAttr(accessor, "enum.0", configVar.Schema.Enum[0]),
						resource.TestCheckResourceAttr(accessor, "enum.1", configVar.Schema.Enum[1]),
					),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ConfigurationVariableCreate(
				client.ConfigurationVariableCreateParams{
					Name:        configVar.Name,
					Value:       configVar.Value,
					IsSensitive: false,
					Scope:       client.ScopeGlobal,
					ScopeId:     "",
					Type:        client.ConfigurationVariableTypeEnvironment,
					EnumValues:  configVar.Schema.Enum,
					Description: configVar.Description,
					Format:      client.Text,
				}).Times(1).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariablesById(configVar.Id).Times(1).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)
		})
	})

	for _, format := range []client.Format{client.HCL, client.JSON} {
		t.Run("Create "+string(format)+" Variable", func(t *testing.T) {

			expectedVariable := `{
A = "A"
B = "B"
C = "C"
}
`

			schema := client.ConfigurationVariableSchema{
				Type:   "string",
				Format: format,
			}
			configVar := client.ConfigurationVariable{
				Id:          "id0",
				Name:        "name0",
				Description: "desc0",
				Value:       expectedVariable,
				Schema:      &schema,
			}
			terraformDirective := `<<EOT
{
%{ for key, value in var.map ~}
${key} = "${value}"
%{ endfor ~}
}
EOT`
			stepConfig := fmt.Sprintf(`
variable "map" {
		description = "a mapped variable"
		type        = map(string)
		default = %s
	}


resource "%s" "test" {
		name = "%s"
		description = "%s"
		value = %s
		format = "%s"
}`, expectedVariable, resourceType, configVar.Name, configVar.Description, terraformDirective, string(format))

			createTestCase := resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: stepConfig,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "value", configVar.Value),
							resource.TestCheckResourceAttr(accessor, "format", string(format)),
						),
					},
				},
			}

			runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
				mock.EXPECT().ConfigurationVariableCreate(
					client.ConfigurationVariableCreateParams{
						Name:        configVar.Name,
						Value:       expectedVariable,
						IsSensitive: false,
						Scope:       client.ScopeGlobal,
						ScopeId:     "",
						Type:        client.ConfigurationVariableTypeEnvironment,
						EnumValues:  configVar.Schema.Enum,
						Description: configVar.Description,
						Format:      format,
					}).Times(1).Return(configVar, nil)
				mock.EXPECT().ConfigurationVariablesById(configVar.Id).Times(1).Return(configVar, nil)
				mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)
			})
		})
	}

	t.Run("Create Enum with wrong value", func(t *testing.T) {
		stepConfig := fmt.Sprintf(`
	resource "%s" "test" {
		name = "%s"
		description = "%s"
		value= "%s"
		enum = ["a","b"]
	}`, resourceType, configVar.Name, configVar.Description, configVar.Value)
		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile(fmt.Sprintf("value - '%s' is not one of the enum options", configVar.Value)),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("Create Enum with empty value", func(t *testing.T) {
		stepConfig := fmt.Sprintf(`
		resource "%s" "test" {
			name = "%s"
			description = "%s"
			value= "%s"
			enum = ["a",""]
		}`, resourceType, configVar.Name, configVar.Description, configVar.Value)

		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile(`an empty enum value is not allowed \(at index 1\)`),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("Create with wrong type", func(t *testing.T) {
		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":  configVar.Name,
						"value": configVar.Value,
						"type":  6,
					}),
					ExpectError: regexp.MustCompile("unknown configuration variable type 6"),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {

		})
	})

	t.Run("Read with wrong api error", func(t *testing.T) {
		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile("could not get configuration variable: error"),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ConfigurationVariableCreate(
				configurationVariableCreateParams).Times(1).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariablesById(configVar.Id).Times(1).Return(client.ConfigurationVariable{}, errors.New("error"))
			mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)
		})
	})

	t.Run("Create api client error", func(t *testing.T) {
		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile(`(could not create configurationVariable: error)`),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ConfigurationVariableCreate(configurationVariableCreateParams).Times(1).Return(client.ConfigurationVariable{}, errors.New("error"))
		})
	})

	t.Run("Configuration Removed in UI", func(t *testing.T) {
		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfig,
				},
				{
					Config: stepConfig,
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ConfigurationVariableCreate(configurationVariableCreateParams).Times(1).Return(configVar, nil),
				mock.EXPECT().ConfigurationVariablesById(configVar.Id).Times(1).Return(configVar, nil),
				mock.EXPECT().ConfigurationVariablesById(configVar.Id).Times(1).Return(client.ConfigurationVariable{}, http.NewMockFailedResponseError(404)),
				mock.EXPECT().ConfigurationVariableCreate(configurationVariableCreateParams).Times(1).Return(configVar, nil),
				mock.EXPECT().ConfigurationVariablesById(configVar.Id).Times(1).Return(configVar, nil),
				mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Update", func(t *testing.T) {
		newIsReadOnly := false
		newIsRequired := true
		newConfigVar := client.ConfigurationVariable{
			Id:          configVar.Id,
			Name:        configVar.Name,
			Description: configVar.Description,
			Value:       "I want to be the config value",
			IsReadOnly:  &newIsReadOnly,
			IsRequired:  &newIsRequired,
		}

		updateTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":         configVar.Name,
						"description":  configVar.Description,
						"value":        configVar.Value,
						"is_read_only": strconv.FormatBool(*configVar.IsReadOnly),
						"is_required":  strconv.FormatBool(*configVar.IsRequired),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", configVar.Id),
						resource.TestCheckResourceAttr(accessor, "description", configVar.Description),
						resource.TestCheckResourceAttr(accessor, "name", configVar.Name),
						resource.TestCheckResourceAttr(accessor, "value", configVar.Value),
						resource.TestCheckResourceAttr(accessor, "is_read_only", strconv.FormatBool(*configVar.IsReadOnly)),
						resource.TestCheckResourceAttr(accessor, "is_required", strconv.FormatBool(*configVar.IsRequired)),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":         newConfigVar.Name,
						"description":  newConfigVar.Description,
						"value":        newConfigVar.Value,
						"format":       client.HCL,
						"is_read_only": strconv.FormatBool(*newConfigVar.IsReadOnly),
						"is_required":  strconv.FormatBool(*newConfigVar.IsRequired),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", newConfigVar.Id),
						resource.TestCheckResourceAttr(accessor, "name", newConfigVar.Name),
						resource.TestCheckResourceAttr(accessor, "description", newConfigVar.Description),
						resource.TestCheckResourceAttr(accessor, "value", newConfigVar.Value),
						resource.TestCheckResourceAttr(accessor, "format", string(client.HCL)),
						resource.TestCheckResourceAttr(accessor, "is_read_only", strconv.FormatBool(*newConfigVar.IsReadOnly)),
						resource.TestCheckResourceAttr(accessor, "is_required", strconv.FormatBool(*newConfigVar.IsRequired)),
					),
				},
			},
		}

		runUnitTest(t, updateTestCase, func(mock *client.MockApiClientInterface) {
			createParams := configurationVariableCreateParams
			updateParams := createParams
			updateParams.Name = newConfigVar.Name
			updateParams.Value = newConfigVar.Value
			updateParams.Description = newConfigVar.Description
			updateParams.Format = client.HCL
			updateParams.IsReadOnly = *newConfigVar.IsReadOnly
			updateParams.IsRequired = *newConfigVar.IsRequired

			mock.EXPECT().ConfigurationVariableCreate(createParams).Times(1).Return(configVar, nil)
			gomock.InOrder(
				mock.EXPECT().ConfigurationVariablesById(configVar.Id).Times(2).Return(configVar, nil),
				mock.EXPECT().ConfigurationVariablesById(configVar.Id).Return(newConfigVar, nil),
			)
			mock.EXPECT().ConfigurationVariableUpdate(client.ConfigurationVariableUpdateParams{CommonParams: updateParams, Id: newConfigVar.Id}).Times(1).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)
		})
	})

	t.Run("Update with wrong type", func(t *testing.T) {
		wrongType := client.ConfigurationVariableType(6)
		newConfigVar := client.ConfigurationVariable{
			Id:    configVar.Id,
			Name:  configVar.Name,
			Value: "I want to be the config value",
			Type:  &wrongType,
		}

		updateTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":         configVar.Name,
						"description":  configVar.Description,
						"value":        configVar.Value,
						"is_read_only": strconv.FormatBool(*configVar.IsReadOnly),
						"is_required":  strconv.FormatBool(*configVar.IsRequired),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", configVar.Id),
						resource.TestCheckResourceAttr(accessor, "name", configVar.Name),
						resource.TestCheckResourceAttr(accessor, "description", configVar.Description),
						resource.TestCheckResourceAttr(accessor, "value", configVar.Value),
						resource.TestCheckResourceAttr(accessor, "is_read_only", strconv.FormatBool(*configVar.IsReadOnly)),
						resource.TestCheckResourceAttr(accessor, "is_required", strconv.FormatBool(*configVar.IsRequired)),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        newConfigVar.Name,
						"description": newConfigVar.Description,
						"value":       newConfigVar.Value,
						"type":        "boooo",
					}),
					ExpectError: regexp.MustCompile("unknown configuration variable type boooo"),
				},
			},
		}

		runUnitTest(t, updateTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ConfigurationVariableCreate(configurationVariableCreateParams).Times(1).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariablesById(configVar.Id).Times(2).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)
		})
	})

	importStateId_id := `{  "Scope": "BLUEPRINT",  "ScopeId": "id0", "Id": "id1", "name": "name0"}`
	importStateId_name := `{  "Scope": "BLUEPRINT",  "ScopeId": "id0",  "name": "name0"}`
	ResourceNameImport := "env0_configuration_variable.test"
	configVarImport := client.ConfigurationVariable{
		Id:          "id1",
		Name:        "name0",
		Description: "desc0",
		Value:       "Variable",
		IsReadOnly:  &isReadonly,
		IsRequired:  &isRequired,
		Scope:       "BLUEPRINT",
	}
	stepConfirImport := resourceConfigCreate(resourceType, resourceName, map[string]interface{}{

		"name":         configVarImport.Name,
		"description":  configVarImport.Description,
		"value":        configVarImport.Value,
		"is_read_only": strconv.FormatBool(*configVar.IsReadOnly),
		"is_required":  strconv.FormatBool(*configVar.IsRequired),
		"template_id":  "id0",
	})

	configurationVariableCreateParamsImport := client.ConfigurationVariableCreateParams{
		Name:        configVarImport.Name,
		Value:       configVarImport.Value,
		IsSensitive: false,
		Scope:       client.ScopeTemplate,
		ScopeId:     "id0",
		Type:        client.ConfigurationVariableTypeEnvironment,
		EnumValues:  nil,
		Description: configVarImport.Description,
		Format:      client.Text,
		IsRequired:  *configVarImport.IsRequired,
		IsReadOnly:  *configVarImport.IsReadOnly,
	}
	t.Run("import by name", func(t *testing.T) {

		createTestCaseForImport := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfirImport,
				},
				{
					ResourceName:            ResourceNameImport,
					ImportState:             true,
					ImportStateId:           importStateId_name,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"is_required", "is_read_only"},
				},
			},
		}

		runUnitTest(t, createTestCaseForImport, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ConfigurationVariableCreate(configurationVariableCreateParamsImport).Times(1).Return(configVarImport, nil)
			mock.EXPECT().ConfigurationVariablesById(configVarImport.Id).Times(2).Return(configVarImport, nil)
			mock.EXPECT().ConfigurationVariablesByScope(client.ScopeTemplate, configurationVariableCreateParamsImport.ScopeId).AnyTimes().Return([]client.ConfigurationVariable{configVarImport}, nil)
			mock.EXPECT().ConfigurationVariableDelete(configVarImport.Id).Times(1).Return(nil)
		})
	})

	t.Run("import by id", func(t *testing.T) {

		createTestCaseForImport := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfirImport,
				},
				{
					ResourceName:            ResourceNameImport,
					ImportState:             true,
					ImportStateId:           importStateId_id,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"is_required", "is_read_only"},
				},
			},
		}

		runUnitTest(t, createTestCaseForImport, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ConfigurationVariableCreate(configurationVariableCreateParamsImport).Times(1).Return(configVarImport, nil)
			mock.EXPECT().ConfigurationVariablesById(configVarImport.Id).Times(3).Return(configVarImport, nil)
			mock.EXPECT().ConfigurationVariableDelete(configVarImport.Id).Times(1).Return(nil)
		})
	})

	t.Run("cant be empty value when isRequired and isReadOnly are true", func(t *testing.T) {
		trueVariable := true
		configVar := client.ConfigurationVariable{
			Id:          "id0",
			Name:        "name0",
			Description: "desc0",
			Value:       "",
			IsReadOnly:  &trueVariable,
			IsRequired:  &trueVariable,
		}
		stepConfig := resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
			"name":         configVar.Name,
			"description":  configVar.Description,
			"value":        configVar.Value,
			"is_read_only": strconv.FormatBool(*configVar.IsReadOnly),
			"is_required":  strconv.FormatBool(*configVar.IsRequired),
		})

		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile("'value' cannot be empty when 'is_read_only' and 'is_required' are true "),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("cant be empty value when isRequired and isReadOnly are true - Update", func(t *testing.T) {
		newIsReadOnly := true
		newIsRequired := true
		newConfigVar := client.ConfigurationVariable{
			Id:          configVar.Id,
			Name:        configVar.Name,
			Description: configVar.Description,
			Value:       configVar.Value,
			IsReadOnly:  &newIsReadOnly,
			IsRequired:  &newIsRequired,
		}

		updateTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":         configVar.Name,
						"description":  configVar.Description,
						"value":        configVar.Value,
						"is_read_only": strconv.FormatBool(newIsReadOnly),
						"is_required":  strconv.FormatBool(newIsRequired),
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", configVar.Id),
						resource.TestCheckResourceAttr(accessor, "description", configVar.Description),
						resource.TestCheckResourceAttr(accessor, "name", configVar.Name),
						resource.TestCheckResourceAttr(accessor, "value", configVar.Value),
						resource.TestCheckResourceAttr(accessor, "is_read_only", strconv.FormatBool(newIsReadOnly)),
						resource.TestCheckResourceAttr(accessor, "is_required", strconv.FormatBool(newIsRequired)),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":         newConfigVar.Name,
						"description":  newConfigVar.Description,
						"value":        "",
						"format":       client.HCL,
						"is_read_only": strconv.FormatBool(newIsReadOnly),
						"is_required":  strconv.FormatBool(newIsRequired),
					}),
					ExpectError: regexp.MustCompile("'value' cannot be empty when 'is_read_only' and 'is_required' are true "),
				},
			},
		}

		runUnitTest(t, updateTestCase, func(mock *client.MockApiClientInterface) {
			createParams := configurationVariableCreateParams
			createParams.IsReadOnly = true
			createParams.IsRequired = true

			mock.EXPECT().ConfigurationVariableCreate(createParams).Times(1).Return(newConfigVar, nil)
			gomock.InOrder(
				mock.EXPECT().ConfigurationVariablesById(newConfigVar.Id).Times(2).Return(newConfigVar, nil),
			)
			mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(1).Return(nil)
		})
	})

	t.Run("When soft delete is on, it should not actually delete", func(t *testing.T) {
		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":         configVar.Name,
						"description":  configVar.Description,
						"value":        configVar.Value,
						"is_read_only": strconv.FormatBool(*configVar.IsReadOnly),
						"is_required":  strconv.FormatBool(*configVar.IsRequired),
						"soft_delete":  true,
					}),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ConfigurationVariableCreate(configurationVariableCreateParams).Times(1).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariablesById(configVar.Id).Times(1).Return(configVar, nil)
			mock.EXPECT().ConfigurationVariableDelete(configVar.Id).Times(0)
		})
	})
}
