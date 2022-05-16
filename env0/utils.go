package env0

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

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
			// 'resource' tag found. Override to tag value.
			fieldNameSC = resFieldName
		}

		field := val.Field(i)
		fieldType := field.Type()

		dval, ok := d.GetOk(fieldNameSC)
		if !ok {
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
			switch fieldType {
			case reflect.TypeOf([]client.ModuleSshKey{}):
				sshKeys := []client.ModuleSshKey{}
				for _, sshKey := range dval.([]interface{}) {
					sshKeys = append(sshKeys, client.ModuleSshKey{
						Name: sshKey.(map[string]interface{})["name"].(string),
						Id:   sshKey.(map[string]interface{})["id"].(string)})
				}
				field.Set(reflect.ValueOf(sshKeys))
			case reflect.TypeOf([]string{}):
				strs := []string{}
				for _, str := range dval.([]interface{}) {
					strs = append(strs, str.(string))
				}
				field.Set(reflect.ValueOf(strs))
			}

		case reflect.String, reflect.Bool, reflect.Int:
			field.Set(reflect.ValueOf(dval).Convert(fieldType))
		default:
			return fmt.Errorf("internal error - unhandled field kind %v", fieldType.Kind())
		}
	}

	return nil
}

// Extracts values from the interface, and writes it to resourcedata.
func writeResourceData(i interface{}, d *schema.ResourceData) error {
	val := reflect.ValueOf(i).Elem()

	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Type().Field(i).Name
		// Assumes golang is CamalCase and Terraform is snake_case.
		// This behavior can be overrided be used in the 'tfschema' tag.
		fieldNameSC := toSnakeCase(fieldName)
		if resFieldName, ok := val.Type().Field(i).Tag.Lookup("tfschema"); ok {
			// 'resource' tag found. Override to tag value.
			fieldNameSC = resFieldName
		}

		field := val.Field(i)
		fieldType := field.Type()

		if fieldName == "Id" {
			id := field.String()
			if len(id) == 0 {
				return errors.New("id is empty")
			}
			d.SetId(id)
			continue
		}

		if d.Get(fieldNameSC) == nil {
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
			if err := d.Set(fieldNameSC, field.String()); err != nil {
				return err
			}
		case reflect.Int:
			if err := d.Set(fieldNameSC, field.Int()); err != nil {
				return err
			}
		case reflect.Bool:
			if err := d.Set(fieldNameSC, field.Bool()); err != nil {
				return err
			}
		case reflect.Slice:
			switch field.Type() {
			case reflect.TypeOf([]client.ModuleSshKey{}):
				var rawSshKeys []map[string]string
				for i := 0; i < field.Len(); i++ {
					sshKey := field.Index(i).Interface().(client.ModuleSshKey)
					rawSshKeys = append(rawSshKeys, map[string]string{"id": sshKey.Id, "name": sshKey.Name})
				}
				if err := d.Set(fieldNameSC, rawSshKeys); err != nil {
					return err
				}
			case reflect.TypeOf([]client.Agent{}):
				var agents []map[string]string
				for i := 0; i < field.Len(); i++ {
					agent := field.Index(i).Interface().(client.Agent)
					agents = append(agents, map[string]string{"agent_key": agent.AgentKey})
				}
				if err := d.Set(fieldNameSC, agents); err != nil {
					return err
				}
			case reflect.TypeOf([]string{}):
				var strs []interface{}
				for i := 0; i < field.Len(); i++ {
					str := field.Index(i).Interface().(string)
					strs = append(strs, str)
				}
				if err := d.Set(fieldNameSC, strs); err != nil {
					return err
				}
			default:
				return fmt.Errorf("internal error - unhandled slice type %v", field.Type())
			}
		default:
			return fmt.Errorf("internal error - unhandled field kind %v", field.Kind())
		}
	}

	return nil
}

func safeSet(d *schema.ResourceData, k string, v interface{}) {
	// Checks that the key exist in the schema before setting the value.
	if test := d.Get(k); test != nil {
		d.Set(k, v)
	}
}
