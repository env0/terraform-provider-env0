package env0

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
var matchTtl = regexp.MustCompile("^([1-9][0-9]*)-([M|w|d|h])$")

type CustomResourceDataField interface {
	ReadResourceData(fieldName string, d *schema.ResourceData) error
	WriteResourceData(fieldName string, d *schema.ResourceData) error
}

type ResourceDataSliceStructValueWriter interface {
	ResourceDataSliceStructValueWrite(map[string]any) error
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

func stringInSlice(str string, strs []string) bool {
	return slices.Contains(strs, str)
}

func reasourceDataGetValue(fieldName string, omitEmpty bool, d *schema.ResourceData) any {
	dval := d.Get(fieldName)
	if dval == nil {
		return nil
	}

	ival, okInt := dval.(int)

	if omitEmpty && okInt && ival == 0 {
		return nil
	}

	sval, okString := dval.(string)
	if omitEmpty && okString && (sval == "") {
		return nil
	}

	_, okBool := dval.(bool)

	if okString || okBool || okInt {
		//nolint:staticcheck // https://github.com/hashicorp/terraform-plugin-sdk/issues/817
		if _, exists := d.GetOkExists(fieldName); !exists {
			return nil
		}
	}

	if s, ok := dval.(*schema.Set); ok && s.Len() == 0 {
		return nil
	}

	if s, ok := dval.([]any); ok && len(s) == 0 {
		return nil
	}

	return dval
}

// Extracts values from the resourcedata, and writes it to the interface.
// Prepends prefix to the fieldName.
func readResourceDataEx(prefix string, i any, d *schema.ResourceData) error {
	// TODO: add a mechanism that returns an error if fields were set in the resourceData but not in the struct.
	// Blocked by: https://github.com/hashicorp/terraform-plugin-sdk/issues/910
	val := reflect.ValueOf(i).Elem()
	for i := range val.NumField() {
		parsedField := getFieldName(val.Type().Field(i))
		if parsedField.skip {
			continue
		}

		fieldName := parsedField.name
		if prefix != "" {
			fieldName = prefix + "." + fieldName
		}

		dval := reasourceDataGetValue(fieldName, parsedField.omitEmpty, d)
		if dval == nil {
			continue
		}

		field := val.Field(i)

		fieldType := field.Type()

		// custom field is a pointer.
		if _, ok := field.Interface().(CustomResourceDataField); ok {
			if field.IsNil() {
				// Init the field a valid value (instead of nil).
				field.Set(reflect.New(field.Type().Elem()))
			}

			if err := field.Interface().(CustomResourceDataField).ReadResourceData(fieldName, d); err != nil {
				return err
			}

			continue
		}

		// custom field is a value, check the pointer.
		if customField, ok := field.Addr().Interface().(CustomResourceDataField); ok {
			if err := customField.ReadResourceData(fieldName, d); err != nil {
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
			var sliceData []any
			if s, ok := dval.(*schema.Set); ok {
				sliceData = s.List()
			} else {
				sliceData = dval.([]any)
			}

			if err := readResourceDataSliceEx(field, sliceData); err != nil {
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

// Extracts values from the resourcedata, and writes it to the interface.
func readResourceData(i any, d *schema.ResourceData) error {
	return readResourceDataEx("", i, d)
}

func readResourceDataSliceStructHelper(field reflect.Value, resource any) error {
	val := field.Elem()
	m := resource.(map[string]any)

	for i := range val.NumField() {
		parsedField := getFieldName(val.Type().Field(i))
		if parsedField.skip {
			continue
		}

		fieldValue, ok := m[parsedField.name]
		if !ok {
			continue
		}

		field := val.Field(i)
		field.Set(reflect.ValueOf(fieldValue))
	}

	return nil
}

// Extracts a list of values from the resourcedata, and writes it to the struct field.
func readResourceDataSliceEx(field reflect.Value, resources []any) error {
	elemType := field.Type().Elem()
	vals := reflect.MakeSlice(field.Type(), 0, len(resources))

	for _, resource := range resources {
		var val reflect.Value

		switch elemType.Kind() {
		case reflect.String:
			val = reflect.ValueOf(resource.(string))
		case reflect.Int:
			val = reflect.ValueOf(resource.(int))
		case reflect.Struct:
			val = reflect.New(elemType)
			if err := readResourceDataSliceStructHelper(val, resource); err != nil {
				return err
			}

			val = val.Elem()
		default:
			return fmt.Errorf("internal error - unhandled slice element kind %v", elemType.Kind())
		}

		vals = reflect.Append(vals, val)
	}

	field.Set(vals)

	return nil
}

type parsedField struct {
	name      string
	skip      bool
	omitEmpty bool
}

// Returns the field name or skip.
func getFieldName(field reflect.StructField) *parsedField {
	var res parsedField

	// Assumes golang is CamalCase and Terraform is snake_case.
	// This behavior can be overridden be used in the 'tfschema' tag.
	res.name = toSnakeCase(field.Name)

	if tag, ok := field.Tag.Lookup("tfschema"); ok {
		if tag == "-" {
			res.skip = true
		} else {
			tagParts := strings.Split(tag, ",")
			nameFromTag := tagParts[0]

			// Override name by tag value.
			if len(nameFromTag) > 0 {
				res.name = nameFromTag
			}

			if len(tagParts) > 1 && tagParts[1] == "omitempty" {
				res.omitEmpty = true
			}
		}
	}

	return &res
}

// Extracts values from the interface, and writes it to resourcedata.
func writeResourceData(i any, d *schema.ResourceData) error {
	val := reflect.ValueOf(i).Elem()

	for i := range val.NumField() {
		parsedField := getFieldName(val.Type().Field(i))
		if parsedField.skip {
			continue
		}

		field := val.Field(i)
		fieldType := field.Type()
		fieldName := parsedField.name

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

		if customField, ok := field.Interface().(CustomResourceDataField); ok && !field.IsNil() {
			if err := customField.WriteResourceData(fieldName, d); err != nil {
				return err
			}

			continue
		}

		if fieldType.Kind() == reflect.Ptr {
			if field.IsNil() {
				if !parsedField.omitEmpty {
					if err := d.Set(fieldName, nil); err != nil {
						return err
					}
				}

				continue
			}

			field = field.Elem()
			fieldType = field.Type()
		}

		switch fieldType.Kind() {
		case reflect.String:
			if parsedField.omitEmpty && field.String() == "" {
				continue
			}

			if err := d.Set(fieldName, field.String()); err != nil {
				return err
			}
		case reflect.Int:
			if parsedField.omitEmpty && field.Int() == 0 {
				continue
			}

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

func getInterfaceSliceValues(i any) []any {
	result := []any{}

	val := reflect.ValueOf(i)

	for i := range val.Len() {
		field := val.Index(i)
		result = append(result, field.Interface())
	}

	return result
}

func writeResourceDataGetSliceValues(i any, name string, d *schema.ResourceData) ([]any, error) {
	ivalues := getInterfaceSliceValues(i)

	var values []any

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
				return nil, err
			}

			values = append(values, value)
		default:
			return nil, fmt.Errorf("internal error - unhandled slice kind %v", valType.Kind())
		}
	}

	return values, nil
}

func getResourceDataSliceStructValue(val reflect.Value, name string, d *schema.ResourceData) (map[string]any, error) {
	value := make(map[string]any)

	for i := range val.NumField() {
		field := val.Field(i)
		fieldType := field.Type()

		if writer, ok := field.Interface().(ResourceDataSliceStructValueWriter); ok {
			if err := writer.ResourceDataSliceStructValueWrite(value); err != nil {
				return nil, err
			}

			continue
		}

		parsedField := getFieldName(val.Type().Field(i))
		if parsedField.skip {
			continue
		}

		fieldName := parsedField.name

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

		if fieldType.Kind() == reflect.Slice {
			values, err := writeResourceDataGetSliceValues(field.Interface(), name+".*."+fieldName, d)
			if err != nil {
				return nil, err
			}

			value[fieldName] = values

			continue
		}

		value[fieldName] = field.Interface()
	}

	return value, nil
}

// Extracts values from a slice of interfaces, and writes it to resourcedata at name.
func writeResourceDataSlice(i any, name string, d *schema.ResourceData) error {
	values, err := writeResourceDataGetSliceValues(i, name, d)
	if err != nil {
		return err
	}

	if values != nil {
		return d.Set(name, values)
	}

	return nil
}

func writeResourceDataEx(prefix string, i any, d *schema.ResourceData) error {
	if prefix == "" {
		return writeResourceData(i, d)
	}

	return writeResourceDataSlice([]any{i}, prefix, d)
}

func ttlToDuration(ttl *string) (time.Duration, error) {
	if ttl == nil || *ttl == "" || *ttl == "Infinite" || *ttl == "inherit" {
		return math.MaxInt64, nil
	}

	match := matchTtl.FindStringSubmatch(*ttl)
	if match == nil {
		return 0, fmt.Errorf("invalid TTL format %s", *ttl)
	}

	numberStr := match[1]

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return 0, fmt.Errorf("invalid TTL format %s - not a number: %w", *ttl, err)
	}
	// M/w/d/h

	hours := number

	switch rangeType := match[2]; rangeType {
	case "M":
		// 'M' varies each month. Assuming it's 30 days.
		hours *= 30 * 24
	case "w":
		hours *= 7 * 24
	case "d":
		hours *= 24
	}

	return time.ParseDuration(fmt.Sprintf("%dh", hours))
}

func lastUnderscoreSplit(s string) []string {
	lastIndex := strings.LastIndex(s, "_")
	if lastIndex == -1 {
		return []string{s}
	}

	return []string{s[:lastIndex], s[lastIndex+1:]}
}
