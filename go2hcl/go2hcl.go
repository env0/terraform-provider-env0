package go2hcl

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// TFSource Enum
type TFSource string

const (
	DataSource TFSource = "data"
	Resource   TFSource = "resource"
)

func DataSourceAccessor(resourceType string, resourceName string) string {
	return hclAccessor(DataSource, resourceType, resourceName)
}

func ResourceAccessor(resourceType string, resourceName string) string {
	return hclAccessor(Resource, resourceType, resourceName)
}

func hclAccessor(source TFSource, resourceType string, resourceName string) string {
	if source == DataSource {
		return fmt.Sprintf("%s.%s.%s", source, resourceType, resourceName)
	}
	return fmt.Sprintf("%s.%s", resourceType, resourceName)
}

func DataSourceConfigCreate(resourceType string, resourceName string, fields map[string]interface{}) string {
	return hclConfigCreate(DataSource, resourceType, resourceName, fields)
}

func ResourceConfigCreate(resourceType string, resourceName string, fields map[string]interface{}) string {
	return hclConfigCreate(Resource, resourceType, resourceName, fields)
}

func hclConfigCreate(source TFSource, resourceType string, resourceName string, fields map[string]interface{}) string {
	hclFields := ""
	for key, value := range fields {
		field, err := toHclField(key, value)
		if err != nil {
			panic(err)
		}
		hclFields += "\n\t" + field
	}
	if hclFields != "" {
		hclFields += "\n"
	}
	return fmt.Sprintf(`%s "%s" "%s" {%s}`, source, resourceType, resourceName, hclFields)
}

func toHclValue(value interface{}) (string, error) {
	reflectedType := reflect.TypeOf(value)
	switch reflectedType.Kind() {
	case reflect.Int:
		return fmt.Sprintf("%v", reflect.ValueOf(value).Interface()), nil
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%v", reflect.ValueOf(value).Interface()), nil
	case reflect.Bool:
		return fmt.Sprintf("%v", reflect.ValueOf(value).Interface()), nil
	case reflect.String:
		return fmt.Sprintf("\"%v\"", reflect.ValueOf(value).Interface()), nil
	case reflect.Slice, reflect.Array:
		listValues := reflect.ValueOf(value)
		var hclValues []string
		for i := 0; i < listValues.Len(); i++ {
			innerValue, err := toHclValue(listValues.Index(i).Interface())
			if err != nil {
				return "", err
			}
			hclValues = append(hclValues, innerValue)
		}
		return fmt.Sprintf("[\n\t%s\n]", strings.Join(hclValues, ",\n\t")), nil
	case reflect.Map:
		mapValue := reflect.ValueOf(value)
		var hclMapFields []string
		for _, key := range mapValue.MapKeys() {
			fieldKey := key.String()
			fieldValue := mapValue.MapIndex(key).Interface()
			hclField, err := toHclField(fieldKey, fieldValue)
			if err != nil {
				return "", err
			}
			hclMapFields = append(hclMapFields, hclField)
		}
		sort.Strings(hclMapFields)
		return fmt.Sprintf("{\n\t%s\n}", strings.Join(hclMapFields, "\n\t")), nil
	default:
		return "", errors.New("can't convert value to hcl")
	}
}

func toHclField(name string, value interface{}) (string, error) {
	hclValue, err := toHclValue(value)
	if err != nil {
		return "", errors.New(fmt.Sprintf("'%s' field has unsupported value - %s", name, err.Error()))
	}
	if reflect.ValueOf(value).Kind() == reflect.Map {
		return fmt.Sprintf("%s %s", name, hclValue), nil
	}
	return fmt.Sprintf("%s = %s", name, hclValue), nil
}
