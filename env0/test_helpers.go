package env0

import (
	"fmt"
)

// TFSource Enum
type TFSource string

const (
	DataResource TFSource = "data"
	Resource     TFSource = "resource"
)

func DataSourceAccessor(resourceType string, resourceName string) string {
	return hclAccessor(DataResource, resourceType, resourceName)
}

func ResourceAccessor(resourceType string, resourceName string) string {
	return hclAccessor(Resource, resourceType, resourceName)
}

func hclAccessor(source TFSource, resourceType string, resourceName string) string {
	if source == DataResource {
		return fmt.Sprintf("%s.%s.%s", source, resourceType, resourceName)
	}
	return fmt.Sprintf("%s.%s", resourceType, resourceName)
}

func DataSourceConfigCreate(resourceType string, resourceName string, fields map[string]string) string {
	return hclConfigCreate(DataResource, resourceType, resourceName, fields)
}

func ResourceConfigCreate(resourceType string, resourceName string, fields map[string]string) string {
	return hclConfigCreate(Resource, resourceType, resourceName, fields)
}

func hclConfigCreate(source TFSource, resourceType string, resourceName string, fields map[string]string) string {
	hclFields := ""
	for key, value := range fields {
		hclFields += fmt.Sprintf(`\n	"%s" = "%s"`, key, value)
	}
	if hclFields != "" {
		hclFields = "\n"
	}
	return fmt.Sprintf(`%s "%s" "%s" {%s}`, source, resourceType, resourceName, hclFields)
}
