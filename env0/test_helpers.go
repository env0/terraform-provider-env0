package env0

import (
	"fmt"
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

func dataSourceConfigCreate(resourceType string, resourceName string, fields map[string]interface{}) string {
	return hclConfigCreate(DataSource, resourceType, resourceName, fields)
}

func resourceConfigCreate(resourceType string, resourceName string, fields map[string]interface{}) string {
	return hclConfigCreate(Resource, resourceType, resourceName, fields)
}

func hclConfigCreate(source TFSource, resourceType string, resourceName string, fields map[string]interface{}) string {
	hclFields := ""
	for key, value := range fields {
		intValue, intOk := value.(int)
		boolValue, boolOk := value.(bool)
		arrayStrValues, arrayStrOk := value.([]string)
		arrayIntValues, arrayIntOk := value.([]int)

		if intOk {
			hclFields += fmt.Sprintf("\n\t%s = %d", key, intValue)
		} else if boolOk {
			hclFields += fmt.Sprintf("\n\t%s = %t", key, boolValue)
		} else if arrayStrOk {
			arrayValueString := ""
			for _, arrayValue := range arrayStrValues {
				arrayValueString += "\"" + arrayValue + "\","
			}
			arrayValueString = arrayValueString[:len(arrayValueString)-1]

			hclFields += fmt.Sprintf("\n\t%s = [%s]", key, arrayValueString)
		} else if arrayIntOk {
			arrayValueString := ""
			for _, arrayValue := range arrayIntValues {
				arrayValueString += fmt.Sprintf("%d,", arrayValue)
			}
			arrayValueString = arrayValueString[:len(arrayValueString)-1]

			hclFields += fmt.Sprintf("\n\t%s = [%s]", key, arrayValueString)
		} else {
			hclFields += fmt.Sprintf("\n\t%s = \"%s\"", key, value)
		}
	}
	if hclFields != "" {
		hclFields += "\n"
	}
	return fmt.Sprintf(`%s "%s" "%s" {%s}`, source, resourceType, resourceName, hclFields)
}

func missingArgumentTestCase(resourceType string, resourceName string, errorResource map[string]interface{}, missingArgumentKey string) resource.TestCase {
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

func missingArgumentTestCaseForCostCred(resourceType string, resourceName string, errorResource map[string]interface{}) resource.TestCase {
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
