package env0

import (
	"fmt"
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

func dataSourceConfigCreate(resourceType string, resourceName string, stringFields map[string]string, intFields map[string]int, boolFields map[string]bool) string {
	return hclConfigCreate(DataSource, resourceType, resourceName, stringFields, intFields, boolFields)
}

func resourceConfigCreate(resourceType string, resourceName string, stringFields map[string]string, intFields map[string]int, boolFields map[string]bool) string {
	return hclConfigCreate(Resource, resourceType, resourceName, stringFields, intFields, boolFields)
}

func hclConfigCreate(source TFSource, resourceType string, resourceName string, stringFields map[string]string, intFields map[string]int, boolFields map[string]bool) string {
	hclFields := ""
	for key, value := range stringFields {
		hclFields += fmt.Sprintf("\n\t%s = \"%s\"", key, value)
	}
	for key, value := range intFields {
		hclFields += fmt.Sprintf("\n\t%s = \"%d\"", key, value)
	}
	for key, value := range boolFields {
		hclFields += fmt.Sprintf("\n\t%s = \"%t\"", key, value)
	}
	if hclFields != "" {
		hclFields += "\n"
	}
	return fmt.Sprintf(`%s "%s" "%s" {%s}`, source, resourceType, resourceName, hclFields)
}
