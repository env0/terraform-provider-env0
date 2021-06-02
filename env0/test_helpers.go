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

func dataSourceConfigCreate(resourceType string, resourceName string, fields map[string]string) string {
	return hclConfigCreate(DataSource, resourceType, resourceName, fields)
}

func resourceConfigCreate(resourceType string, resourceName string, fields map[string]string) string {
	return hclConfigCreate(Resource, resourceType, resourceName, fields)
}

func hclConfigCreate(source TFSource, resourceType string, resourceName string, fields map[string]string) string {
	hclFields := ""
	for key, value := range fields {
		hclFields += fmt.Sprintf("\n\t\"%s\" = \"%s\"", key, value)
	}
	if hclFields != "" {
		hclFields += "\n"
	}
	return fmt.Sprintf(`%s "%s" "%s" {%s}`, source, resourceType, resourceName, hclFields)
}
