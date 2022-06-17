package env0

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

type CustomResourceDataField interface {
	ReadResourceData(fieldName string, d *schema.ResourceData) error
	WriteResourceData(fieldName string, d *schema.ResourceData) error
}

type ResourceDataSliceStructValueWriter interface {
	ResourceDataSliceStructValueWrite(map[string]interface{}) error
}

// https://stackoverflow.com/questions/56616196/how-to-convert-camel-case-string-to-snake-case
func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

// Extracts values from the resourcedata, and writes it to the interface.
func readResourceData(i interface{}, d *schema.ResourceData) error {
	// TODO: add a mechanism that returns an error if fields were set in the resourceData but not in the struct.
	// Blocked by: https://github.com/hashicorp/terraform-plugin-sdk/issues/910

	val := reflect.ValueOf(i).Elem()
	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Type().Field(i).Name
		// Assumes golang is CamalCase and Terraform is snake_case.
		// This behavior can be overrided be used in the 'tfschema' tag.
		fieldNameSC := toSnakeCase(fieldName)
		if resFieldName, ok := val.Type().Field(i).Tag.Lookup("tfschema"); ok {
			if resFieldName == "-" {
				continue
			}

			// 'resource' tag found. Override to tag value.
			fieldNameSC = resFieldName
		}

		field := val.Field(i)

		fieldType := field.Type()

		dval, ok := d.GetOk(fieldNameSC)
		if !ok {
			continue
		}

		// custom field is a pointer.
		if _, ok := field.Interface().(CustomResourceDataField); ok {
			if field.IsNil() {
				// Init the field a valid value (instead of nil).
				field.Set(reflect.New(field.Type().Elem()))
			}
			if err := field.Interface().(CustomResourceDataField).ReadResourceData(fieldNameSC, d); err != nil {
				return err
			}
			continue
		}

		// custom field is a value, check the pointer.
		if customField, ok := field.Addr().Interface().(CustomResourceDataField); ok {
			if err := customField.ReadResourceData(fieldNameSC, d); err != nil {
				return err
			}
			continue
		}

		switch fieldType.Kind() {
		case reflect.Ptr:
			switch fieldType.Elem().Kind() {
			case reflect.Int:
				i := dval.(int)
				field.Set(reflect.ValueOf(&i))
			case reflect.Bool:
				b := dval.(bool)
				field.Set(reflect.ValueOf(&b))
			case reflect.String:
				s := dval.(string)
				field.Set(reflect.ValueOf(&s))
			default:
				return fmt.Errorf("internal error - unhandled field pointer kind %v", fieldType.Elem().Kind())
			}
		case reflect.Slice:
			if err := readResourceDataSlice(field, dval.([]interface{})); err != nil {
				return err
			}
		case reflect.String, reflect.Bool, reflect.Int:
			field.Set(reflect.ValueOf(dval).Convert(fieldType))
		default:
			return fmt.Errorf("internal error - unhandled field kind %v", fieldType.Kind())
		}
	}

	return nil
}

func readResourceDataSliceStructHelper(field reflect.Value, resource interface{}) error {
	val := field.Elem()
	m := resource.(map[string]interface{})

	for i := 0; i < val.NumField(); i++ {
		fieldName, skip := getFieldName(val.Type().Field(i))
		if skip {
			continue
		}

		fieldValue, ok := m[fieldName]
		if !ok {
			continue
		}

		field := val.Field(i)
		field.Set(reflect.ValueOf(fieldValue))
	}

	return nil
}

// Extracts a list of values from the resourcedata, and writes it to the struct field.
func readResourceDataSlice(field reflect.Value, resources []interface{}) error {
	elemType := field.Type().Elem()
	vals := reflect.MakeSlice(field.Type(), 0, len(resources))

	for _, resource := range resources {
		var val reflect.Value
		switch elemType.Kind() {
		case reflect.String:
			val = reflect.ValueOf(resource.(string))
		case reflect.Struct:
			val = reflect.New(elemType)
			if err := readResourceDataSliceStructHelper(val, resource); err != nil {
				return err
			}
		default:
			return fmt.Errorf("internal error - unhandled slice element kind %v", elemType.Kind())
		}
		vals = reflect.Append(vals, val.Elem())
	}

	field.Set(vals)

	return nil
}

// Returns the field name or skip.
func getFieldName(field reflect.StructField) (string, bool) {
	name := field.Name
	// Assumes golang is CamalCase and Terraform is snake_case.
	// This behavior can be overrided be used in the 'tfschema' tag.
	name = toSnakeCase(name)
	if tag, ok := field.Tag.Lookup("tfschema"); ok {
		if tag == "-" {
			return "", true
		}

		// 'resource' tag found. Override to tag value.
		name = tag
	}

	return name, false
}

// Extracts values from the interface, and writes it to resourcedata.
func writeResourceData(i interface{}, d *schema.ResourceData) error {
	val := reflect.ValueOf(i).Elem()

	for i := 0; i < val.NumField(); i++ {
		fieldName, skip := getFieldName(val.Type().Field(i))
		if skip {
			continue
		}

		field := val.Field(i)
		fieldType := field.Type()

		if fieldName == "id" {
			id := field.String()
			if len(id) == 0 {
				return errors.New("id is empty")
			}
			d.SetId(id)
			continue
		}

		if d.Get(fieldName) == nil {
			continue
		}

		if customField, ok := field.Interface().(CustomResourceDataField); ok {
			if err := customField.WriteResourceData(fieldName, d); err != nil {
				return err
			}
			continue
		}

		if fieldType.Kind() == reflect.Ptr {
			if field.IsNil() {
				continue
			}

			field = field.Elem()
			fieldType = field.Type()
		}

		switch fieldType.Kind() {
		case reflect.String:
			if err := d.Set(fieldName, field.String()); err != nil {
				return err
			}
		case reflect.Int:
			if err := d.Set(fieldName, field.Int()); err != nil {
				return err
			}
		case reflect.Bool:
			if err := d.Set(fieldName, field.Bool()); err != nil {
				return err
			}
		case reflect.Slice:
			if err := writeResourceDataSlice(field.Interface(), fieldName, d); err != nil {
				return err
			}
		default:
			return fmt.Errorf("internal error - unhandled field kind %v", field.Kind())
		}
	}

	return nil
}

func getInterfaceSliceValues(i interface{}) []interface{} {
	var result []interface{}

	val := reflect.ValueOf(i)

	for i := 0; i < val.Len(); i++ {
		field := val.Index(i)
		result = append(result, field.Interface())
	}

	return result
}

func getResourceDataSliceStructValue(val reflect.Value, name string, d *schema.ResourceData) (map[string]interface{}, error) {
	value := make(map[string]interface{})

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := field.Type()

		if writer, ok := field.Interface().(ResourceDataSliceStructValueWriter); ok {
			if err := writer.ResourceDataSliceStructValueWrite(value); err != nil {
				return nil, err
			}
			continue
		}

		fieldName, skip := getFieldName(val.Type().Field(i))
		if skip {
			continue
		}

		if fieldType.Kind() == reflect.Ptr {
			if field.IsNil() {
				continue
			}

			field = field.Elem()
		}

		// Check if the field exist in the schema. `*` is for any index.
		if d.Get(name+".*."+fieldName) == nil {
			continue
		}

		value[fieldName] = field.Interface()
	}

	return value, nil
}

// Extracts values from a slice of interfaces, and writes it to resourcedata at name.
func writeResourceDataSlice(i interface{}, name string, d *schema.ResourceData) error {
	ivalues := getInterfaceSliceValues(i)
	var values []interface{}

	// Iterate over the slice of values and build a slice of terraform values.
	for _, ivalue := range ivalues {
		val := reflect.ValueOf(ivalue)
		valType := val.Type()

		if valType.Kind() == reflect.Ptr {
			if val.IsNil() {
				continue
			}

			val = val.Elem()
			valType = val.Type()
		}

		switch valType.Kind() {
		case reflect.String:
			values = append(values, val.String())
		case reflect.Int:
			values = append(values, val.Int())
		case reflect.Bool:
			values = append(values, val.Bool())
		case reflect.Struct:
			// Slice of structs.
			value, err := getResourceDataSliceStructValue(val, name, d)
			if err != nil {
				return err
			}
			values = append(values, value)
		default:
			return fmt.Errorf("internal error - unhandled slice kind %v", valType.Kind())
		}
	}

	if values != nil {
		return d.Set(name, values)
	}

	return nil
}

func safeSet(d *schema.ResourceData, k string, v interface{}) {
	// Checks that the key exist in the schema before setting the value.
	if test := d.Get(k); test != nil {
		d.Set(k, v)
	}
}
