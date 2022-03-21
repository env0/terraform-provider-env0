package env0

import (
	"errors"
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
	val := reflect.ValueOf(i).Elem()
	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Type().Field(i).Name
		// Assumes golang is CamalCase and Terraform is snake_case.
		fieldNameSC := toSnakeCase(fieldName)
		field := val.Field(i)
		fieldType := field.Type()

		dval, ok := d.GetOk(fieldNameSC)
		if !ok {
			continue
		}

		switch fieldType.Kind() {
		case reflect.Ptr:
			switch field.Type().Elem().Kind() {
			case reflect.Int:
				i := dval.(int)
				field.Set(reflect.ValueOf(&i))
			case reflect.Bool:
				b := dval.(bool)
				field.Set(reflect.ValueOf(&b))
			}
		case reflect.Slice:
			switch field.Type() {
			case reflect.TypeOf([]client.ModuleSshKey{}):
				sshKeys := []client.ModuleSshKey{}
				for _, sshKey := range dval.([]interface{}) {
					sshKeys = append(sshKeys, client.ModuleSshKey{
						Name: sshKey.(map[string]interface{})["name"].(string),
						Id:   sshKey.(map[string]interface{})["id"].(string)})
				}
				field.Set(reflect.ValueOf(sshKeys))
			}
		default:
			field.Set(reflect.ValueOf(dval))
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
		fieldNameSC := toSnakeCase(fieldName)
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
			}
		}
	}

	return nil
}
