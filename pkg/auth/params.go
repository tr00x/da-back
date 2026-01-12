package auth

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func BuildParams(v any) (keys []string, values []string, args []any) {
	placeHolder := 1

	if v == nil {
		return nil, nil, nil
	}

	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, nil, nil
		}
		val = val.Elem()
		typ = typ.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, nil, nil
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i).Interface()

		if value == nil {
			continue
		}

		if reflect.ValueOf(value).IsZero() {
			continue
		}

		jsonTag := field.Tag.Get("json")

		if jsonTag == "" {
			return nil, nil, nil
		}

		// Handle arrays for PostgreSQL
		if reflect.TypeOf(value).Kind() == reflect.Slice {
			// Skip empty slices
			sliceValue := reflect.ValueOf(value)
			if sliceValue.Len() == 0 {
				continue
			}
		}

		keys = append(keys, jsonTag)
		values = append(values, "$"+strconv.Itoa(placeHolder))
		placeHolder++
		args = append(args, value)
	}
	return keys, values, args
}

func QueryParamToArray(param string) []string {
	if param == "" {
		return nil
	}
	parts := strings.Split(param, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func QueryParamToIntArray(param string) ([]int, error) {
	if param == "" {
		return nil, nil
	}

	parts := strings.Split(param, ",")
	result := make([]int, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		result = append(result, n)
	}
	return result, nil
}

// ToPostgreSQLArray converts a Go slice to a format suitable for PostgreSQL array insertion
// This is useful when you need to manually construct array values for SQL queries
func ToPostgreSQLArray(slice any) string {
	if slice == nil {
		return "NULL"
	}

	val := reflect.ValueOf(slice)
	if val.Kind() != reflect.Slice {
		return "NULL"
	}

	if val.Len() == 0 {
		return "ARRAY[]"
	}

	var elements []string
	for i := 0; i < val.Len(); i++ {
		element := val.Index(i).Interface()
		switch v := element.(type) {
		case string:
			elements = append(elements, "'"+strings.ReplaceAll(v, "'", "''")+"'")
		case int, int32, int64:
			elements = append(elements, strconv.FormatInt(reflect.ValueOf(v).Int(), 10))
		case float32, float64:
			elements = append(elements, strconv.FormatFloat(reflect.ValueOf(v).Float(), 'f', -1, 64))
		default:
			elements = append(elements, "'"+strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''")+"'")
		}
	}

	return "ARRAY[" + strings.Join(elements, ", ") + "]"
}
