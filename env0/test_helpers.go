package env0

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TFSource Enum
type TFSource string

const (
	DataSource TFSource = "data"
	Resource   TFSource = "resource"
)

func dataSourceAccessor(resourceType string, resourceName string) string {
	return hclAccessor(DataSource, resourceType, resourceName)
}

func resourceAccessor(resourceType string, resourceName string) string {
	return hclAccessor(Resource, resourceType, resourceName)
}

func hclAccessor(source TFSource, resourceType string, resourceName string) string {
	if source == DataSource {
		return fmt.Sprintf("%s.%s.%s", source, resourceType, resourceName)
	}

	return fmt.Sprintf("%s.%s", resourceType, resourceName)
}

func dataSourceConfigCreate(resourceType string, resourceName string, fields map[string]any) string {
	return hclConfigCreate(DataSource, resourceType, resourceName, fields)
}

func resourceConfigCreate(resourceType string, resourceName string, fields map[string]any) string {
	return hclConfigCreate(Resource, resourceType, resourceName, fields)
}

func hclConfigCreate(source TFSource, resourceType string, resourceName string, fields map[string]any) string {
	hclFields := ""

	for key, value := range fields {
		valueType := reflect.TypeOf(value)

		switch valueType {
		case reflect.TypeFor[int]():
			hclFields += fmt.Sprintf("\n\t%s = %d", key, value.(int))
		case reflect.TypeFor[bool]():
			hclFields += fmt.Sprintf("\n\t%s = %t", key, value.(bool))
		case reflect.TypeFor[[]string]():
			arrayValueString := ""

			for _, arrayValue := range value.([]string) {
				arrayValueString += "\"" + arrayValue + "\","
			}

			arrayValueString = arrayValueString[:len(arrayValueString)-1]

			hclFields += fmt.Sprintf("\n\t%s = [%s]", key, arrayValueString)
		case reflect.TypeFor[[]int]():
			arrayValueString := ""

			for _, arrayValue := range value.([]int) {
				arrayValueString += fmt.Sprintf("%d,", arrayValue)
			}

			arrayValueString = arrayValueString[:len(arrayValueString)-1]

			hclFields += fmt.Sprintf("\n\t%s = [%s]", key, arrayValueString)
		default:
			hclFields += fmt.Sprintf("\n\t%s = \"%s\"", key, value)
		}
	}

	if hclFields != "" {
		hclFields += "\n"
	}

	return fmt.Sprintf(`%s "%s" "%s" {%s}`, source, resourceType, resourceName, hclFields)
}

func missingArgumentTestCase(resourceType string, resourceName string, errorResource map[string]any, missingArgumentKey string) resource.TestCase {
	testCaseFormMissingValidInputError := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config:      resourceConfigCreate(resourceType, resourceName, errorResource),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`The argument \"%s\" is required, but no definition was found`, missingArgumentKey)),
			},
		},
	}

	return testCaseFormMissingValidInputError
}

func missingArgumentTestCaseForCostCred(resourceType string, resourceName string, errorResource map[string]any) resource.TestCase {
	testCaseFormMissingValidInputError := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config:      resourceConfigCreate(resourceType, resourceName, errorResource),
				ExpectError: regexp.MustCompile("Error: Missing required argument"),
			},
		},
	}

	return testCaseFormMissingValidInputError
}
