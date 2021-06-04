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
	var hclFields []string
	for key, value := range fields {
		field, err := toHclField(key, value, 1)
		if err != nil {
			panic(err)
		}
		hclFields = append(hclFields, field)
	}
	sort.Strings(hclFields)
	object := "{}"
	if len(hclFields) != 0 {
		object = fmt.Sprintf("{\n%s\n}", strings.Join(hclFields, "\n"))
	}
	return fmt.Sprintf(`%s "%s" "%s" %s`, source, resourceType, resourceName, object)
}

func toHclField(name string, value interface{}, indentLevel uint) (string, error) {
	hclValue, err := toHclValue(value, indentLevel)
	if err != nil {
		return "", errors.New(fmt.Sprintf("'%s' field has unsupported value - %s", name, err.Error()))
	}

	indentation := strings.Repeat("\t", int(indentLevel))
	if reflect.ValueOf(value).Kind() == reflect.Map {
		return fmt.Sprintf("%s%s %s", indentation, name, hclValue), nil
	}
	return fmt.Sprintf("%s%s = %s", indentation, name, hclValue), nil
}

func toHclValue(value interface{}, indentLevel uint) (string, error) {
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
			innerValue, err := toHclValue(listValues.Index(i).Interface(), 0)
			if err != nil {
				return "", err
			}
			hclValues = append(hclValues, innerValue)
		}
		arrayContent := strings.Join(hclValues, ", ")
		return fmt.Sprintf("[%s]", arrayContent), nil
	case reflect.Map:
		mapValue := reflect.ValueOf(value)
		var hclMapFields []string
		for _, key := range mapValue.MapKeys() {
			fieldKey := key.String()
			fieldValue := mapValue.MapIndex(key).Interface()
			hclField, err := toHclField(fieldKey, fieldValue, indentLevel+1)
			if err != nil {
				return "", err
			}
			hclMapFields = append(hclMapFields, hclField)
		}
		sort.Strings(hclMapFields)
		mapContent := strings.Join(hclMapFields, "\n")
		indentation := strings.Repeat("\t", int(indentLevel))
		return fmt.Sprintf("{\n%s\n%s}", mapContent, indentation), nil
	default:
		return "", errors.New("can't convert value to hcl")
	}
}
