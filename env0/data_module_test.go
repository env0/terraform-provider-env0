package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestModuleDataSource(t *testing.T) {
	module := client.Module{
		Id:             "id0",
		ModuleName:     "module0",
		ModuleProvider: "provider0",
		TokenId:        "t0",
		TokenName:      "n0",
		Repository:     "r0",
	}

	otherModule := client.Module{
		Id:             "id1",
		ModuleName:     "module1",
		ModuleProvider: "provider1",
		TokenId:        "t1",
		TokenName:      "n1",
		Repository:     "r1",
	}

	moduleFieldsByName := map[string]interface{}{"module_name": module.ModuleName}
	moduleFieldsById := map[string]interface{}{"id": module.Id}

	resourceType := "env0_module"
	resourceName := "test_module"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]interface{}) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", module.Id),
						resource.TestCheckResourceAttr(accessor, "module_name", module.ModuleName),
						resource.TestCheckResourceAttr(accessor, "module_provider", module.ModuleProvider),
						resource.TestCheckResourceAttr(accessor, "token_id", module.TokenId),
						resource.TestCheckResourceAttr(accessor, "token_name", module.TokenName),
						resource.TestCheckResourceAttr(accessor, "repository", module.Repository),
					),
				},
			},
		}
	}

	getErrorTestCase := func(input map[string]interface{}, expectedError string) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      dataSourceConfigCreate(resourceType, resourceName, input),
					ExpectError: regexp.MustCompile(expectedError),
				},
			},
		}
	}

	mockListModulesCall := func(returnValue []client.Module) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Modules().AnyTimes().Return(returnValue, nil)
		}
	}

	mockModuleCall := func(returnValue *client.Module) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Module(module.Id).AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("By ID", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(moduleFieldsById),
			mockModuleCall(&module),
		)
	})

	t.Run("By Name", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(moduleFieldsByName),
			mockListModulesCall([]client.Module{module, otherModule}),
		)
	})

	t.Run("Throw error when no name or id is supplied", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(map[string]interface{}{}, "one of `id,module_name` must be specified"),
			func(mock *client.MockApiClientInterface) {},
		)
	})

	t.Run("Throw error when by name and more than one module exists", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(moduleFieldsByName, "found multiple modules"),
			mockListModulesCall([]client.Module{module, otherModule, module}),
		)
	})

	t.Run("Throw error when by name and no module found with that name", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(moduleFieldsByName, "not found"),
			mockListModulesCall([]client.Module{otherModule}),
		)
	})
}
