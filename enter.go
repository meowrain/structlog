// Package structlog provides utilities for logging the fields of a struct
// using reflection. It supports nested structs, pointers, and custom tags
// to customize the logged field names.
package structlog

import (
	"fmt"
	"reflect"
	"strings"
)

// LogStructFields parses the fields of a struct and returns a map of field values
// keyed by their "testlog" tag or field name if the tag is not present.
// It handles nested structs and pointers, and skips unexported fields.
//
// Parameters:
//   - v: The struct or pointer to a struct to be parsed.
//
// Returns:
//   - A map[string]interface{} where keys are field names (or tags) and values
//     are the corresponding field values.
func LogStructFields(v interface{}) map[string]interface{} {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	result := make(map[string]interface{})
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Get the "testlog" tag
		tag := field.Tag.Get("testlog")
		if tag == "" {
			tag = field.Name // Use field name if "testlog" tag is not present
		}

		// Skip unexported fields
		if !fieldValue.CanInterface() {
			continue
		}

		// Handle nested structs
		if fieldValue.Kind() == reflect.Struct {
			nested := LogStructFields(fieldValue.Interface())
			for k, v := range nested {
				result[tag+"."+k] = v
			}
			continue
		}

		// Handle pointer types
		if fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				result[tag] = nil
				continue
			}
			fieldValue = fieldValue.Elem()
		}

		result[tag] = fieldValue.Interface()
	}

	return result
}

// LogStruct formats and returns the contents of a struct as a string.
// It uses the "testlog" tag as the key for each field in the output.
//
// Parameters:
//   - v: The struct or pointer to a struct to be logged.
//
// Returns:
//   - A string representation of the struct's fields and their values.
func LogStruct(v interface{}) string {
	fields := LogStructFields(v)
	var sb strings.Builder

	for k, v := range fields {
		sb.WriteString(k)
		sb.WriteString(": ")
		sb.WriteString(fmt.Sprintf("%v", v))
		sb.WriteString("\n")
	}

	return sb.String()
}
