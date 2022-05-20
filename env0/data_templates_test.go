package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestTemplatesDataSource(t *testing.T) {
	template1 := client.Template{
		Id:   "id0",
		Name: "name0",
	}

	template2 := client.Template{
		Id:   "id1",
		Name: "name1",
	}

	template3 := client.Template{
		Id:        "id2",
		Name:      "name2",
		IsDeleted: true,
	}

	resourceType := "env0_templates"
	resourceName := "test_templates"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getTestCase := func() resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "names.0", template1.Name),
						resource.TestCheckResourceAttr(accessor, "names.1", template2.Name),
						resource.TestCheckNoResourceAttr(accessor, "names.2"),
					),
				},
			},
		}
	}

	mockTemplates := func(returnValue []client.Template) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Templates().AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("Success", func(t *testing.T) {
		runUnitTest(t,
			getTestCase(),
			mockTemplates([]client.Template{template1, template2, template3}),
		)
	})

	t.Run("API Call Error", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config:      dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{}),
						ExpectError: regexp.MustCompile("error"),
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Templates().AnyTimes().Return(nil, errors.New("error"))
			},
		)
	})
}
