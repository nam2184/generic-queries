
package util

import (
	"fmt"
	"reflect"
	"time"
)

func GetArgs[T any](item T) ([]interface{}, error) {
    v := reflect.ValueOf(item)
    if v.Kind() != reflect.Struct {
        return nil, fmt.Errorf("expected a struct, got %s", v.Kind())
    }

    var args []interface{} // Use a dynamically sized slice
    for i := 0; i < v.NumField(); i++ {
        value := v.Field(i)

        // Add value to args if it's not a zero value or is explicitly zero (like int 0)
        if !isZeroValue(value) || (value.Kind() == reflect.Int && value.Int() == 0) || (value.Kind() == reflect.Bool && value.Bool() == false) {
            args = append(args, value.Interface())
        }
    }
    return args, nil
}

func PrintStructAttributes(s interface{}) {
	// Dereference pointer if needed
	v := reflect.Indirect(reflect.ValueOf(s))

	// Ensure the underlying value is a struct
	if v.Kind() != reflect.Struct {
		fmt.Println("Provided value is not a struct or a pointer to a struct")
		return
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Skip unexported fields
		if !value.CanInterface() {
			continue
		}

		fmt.Printf("%s: %v\n", field.Name, value.Interface())
	}
}

func IsZeroField(v interface{}, fieldName string) bool {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return false
	}

	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return false	
  }

	return reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface())
}

func GetZero[T any]() T {
    var result T
    return result
}

func IsZero[T any](value T) bool {
    return reflect.DeepEqual(value, reflect.Zero(reflect.TypeOf(value)).Interface())
}

func Fields[T any](entity T) (string, error) {
    // Get the value and type of the struct
    v := reflect.ValueOf(entity)
    t := reflect.TypeOf(entity)

    // Ensure we're working with a struct
    if t.Kind() != reflect.Struct {
        return "", fmt.Errorf("GenerateNamedParams expects a struct, got %s", t.Kind())
    }

    var placeholders []string
    var params []string

    // Iterate over struct fields
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
        dbTag := field.Tag.Get("db")
        
        // Skip unexported fields or fields without a db tag
        if !field.IsExported() || dbTag == "" {
            continue
        }
        
        switch value.Kind() {
        case reflect.Bool:
            placeholders = append(placeholders, fmt.Sprintf(":%s", dbTag))
            params = append(params, dbTag)
        default:
            if !isZeroValue(value) || (value.Kind() == reflect.Int && value.Int() == 0) {
              placeholders = append(placeholders, fmt.Sprintf(":%s", dbTag))
              params = append(params, dbTag)
            }
        }
    }

    // Construct the SQL parameter placeholder string
    placeholderString := fmt.Sprintf("%s", join(params, ", "))
    return placeholderString, nil
}

func Fields2[T any](entity T) (string, error) {
    // Get the value and type of the struct
    v := reflect.ValueOf(entity)
    t := reflect.TypeOf(entity)

    // Ensure we're working with a struct
    if t.Kind() != reflect.Struct {
        return "", fmt.Errorf("GenerateNamedParams expects a struct, got %s", t.Kind())
    }

    var placeholders []string
    var params []string

    // Iterate over struct fields
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
        dbTag := field.Tag.Get("db")
        
        // Skip unexported fields or fields without a db tag
        if !field.IsExported() || dbTag == "" {
            continue
        }
        
        if !isZeroValue(value) {
              placeholders = append(placeholders, fmt.Sprintf(":%s", dbTag))
              params = append(params, dbTag)
        }
    }

    // Construct the SQL parameter placeholder string
    placeholderString := fmt.Sprintf("%s", join(params, ", "))
    return placeholderString, nil
}



func AllFields[T any](entity T) (string, error) {
    // Get the value and type of the struct
    t := reflect.TypeOf(entity)

    // Ensure we're working with a struct
    if t.Kind() != reflect.Struct {
        return "", fmt.Errorf("GenerateNamedParams expects a struct, got %s", t.Kind())
    }

    var placeholders []string
    var params []string

    // Iterate over struct fields
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        dbTag := field.Tag.Get("db")

        
        placeholders = append(placeholders, fmt.Sprintf(":%s", dbTag))
        params = append(params, dbTag)
    }

    // Construct the SQL parameter placeholder string
    placeholderString := fmt.Sprintf("%s", join(params, ", "))
    return placeholderString, nil
}

func FieldsArray[T any](entity T) ([]string, error) {
    t := reflect.TypeOf(entity)

    if t.Kind() != reflect.Struct {
        return nil, fmt.Errorf("GenerateNamedParams expects a struct, got %s", t.Kind())
    }

    var fields []string

    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        dbTag := field.Tag.Get("db")
        
        fields = append(fields, dbTag)
    }

    return fields, nil
}

func GeneratePositionalParams[T any](entity T) (string, error) {
    // Get the value and type of the struct
    v := reflect.ValueOf(entity)
    t := reflect.TypeOf(entity)

    // Ensure we're working with a struct
    if t.Kind() != reflect.Struct {
        return "", fmt.Errorf("GeneratePositionalParams expects a struct, got %s", t.Kind())
    }

    var placeholders []string
    paramIndex := 1 // Start indexing from 1 for PostgreSQL positional placeholders

    // Iterate over struct fields
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
        dbTag := field.Tag.Get("db")

        // Skip unexported fields or fields without a db tag
        if !field.IsExported() || dbTag == "" {
            continue
        }

        // Check if the field should be included in the query
        switch value.Kind() {
        case reflect.Bool:
            placeholders = append(placeholders, fmt.Sprintf("$%d", paramIndex))
            paramIndex++
        default:
            if !isZeroValue(value) || (value.Kind() == reflect.Int && value.Int() == 0) {
                placeholders = append(placeholders, fmt.Sprintf("$%d", paramIndex))
                paramIndex++
            }
        }
    }

    // Join placeholders into a single string
    placeholderString := join(placeholders, ", ")
    return placeholderString, nil
}

func GenerateNamedParams[T any](entity T) (string, error) {
    // Get the value and type of the struct
    v := reflect.ValueOf(entity)
    t := reflect.TypeOf(entity)

    // Ensure we're working with a struct
    if t.Kind() != reflect.Struct {
        return "", fmt.Errorf("GenerateNamedParams expects a struct, got %s", t.Kind())
    }

    var placeholders []string
    var params []string

    // Iterate over struct fields
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
        dbTag := field.Tag.Get("db")

        // Skip unexported fields or fields without a db tag
        if !field.IsExported() || dbTag == "" {
            continue
        }
        
       
        switch value.Kind() {
        case reflect.Bool:
            placeholders = append(placeholders, fmt.Sprintf(":%s", dbTag))
            params = append(params, dbTag)
        default:
            if !isZeroValue(value) || (value.Kind() == reflect.Int && value.Int() == 0) {
              placeholders = append(placeholders, fmt.Sprintf(":%s", dbTag))
              params = append(params, dbTag)
            }
        }
    }

    // Construct the SQL parameter placeholder string
    placeholderString := fmt.Sprintf(":%s", join(params, ", :"))
    return placeholderString, nil
}

func GenerateNamedParams2[T any](entity T) (string, error) {
    // Get the value and type of the struct
    v := reflect.ValueOf(entity)
    t := reflect.TypeOf(entity)

    // Ensure we're working with a struct
    if t.Kind() != reflect.Struct {
        return "", fmt.Errorf("GenerateNamedParams expects a struct, got %s", t.Kind())
    }

    var placeholders []string
    var params []string

    // Iterate over struct fields
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
        dbTag := field.Tag.Get("db")

        // Skip unexported fields or fields without a db tag
        if !field.IsExported() || dbTag == "" {
            continue
        }
        
       
        if !isZeroValue(value) {
              placeholders = append(placeholders, fmt.Sprintf(":%s", dbTag))
              params = append(params, dbTag)
        }
    }

    // Construct the SQL parameter placeholder string
    placeholderString := fmt.Sprintf(":%s", join(params, ", :"))
    return placeholderString, nil
}

func FieldsAndParams[T any](entity T) (string, error) {
    // Get the value and type of the struct
    v := reflect.ValueOf(entity)
    t := reflect.TypeOf(entity)

    // Ensure we're working with a struct
    if t.Kind() != reflect.Struct {
        return "", fmt.Errorf("FieldsAndParams expects a struct, got %s", t.Kind())
    }

    var placeholders []string

    // Iterate over struct fields
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
        dbTag := field.Tag.Get("db")

        // Skip unexported fields or fields without a db tag
        if !field.IsExported() || dbTag == "" {
            continue
        }
        
        // Check if the field has a zero value and decide whether to include it
        if !isZeroValue(value) || (value.Kind() == reflect.Int && value.Int() == 0) {
            // Construct the SQL clause like: "column_name = :column_name"
            placeholders = append(placeholders, fmt.Sprintf("%s = :%s", dbTag, dbTag))
        }
    }

    // Join all placeholders with ", " for the SQL SET clause
    placeholderString := join(placeholders, ", ")
    return placeholderString, nil
}

func FieldsAndPlaceholders[T any](entity T) (string, error) {
    // Get the value and type of the struct
    v := reflect.ValueOf(entity)
    t := reflect.TypeOf(entity)

    // Ensure we're working with a struct
    if t.Kind() != reflect.Struct {
        return "", fmt.Errorf("FieldsAndParams expects a struct, got %s", t.Kind())
    }

    var placeholders []string
    count := 0
    // Iterate over struct fields
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
        dbTag := field.Tag.Get("db")
        // Skip unexported fields or fields without a db tag
        if !field.IsExported() || dbTag == "" {
            continue
        }
        // Check if the field has a zero value and decide whether to include it
        if !isZeroValue(value) || (value.Kind() == reflect.Int && value.Int() == 0) {
            // Construct the SQL clause like: "column_name = :column_name"
            count++  
            placeholders = append(placeholders, fmt.Sprintf("%s = $%d", dbTag, count))
        }
    }

    // Join all placeholders with ", " for the SQL SET clause
    placeholderString := join(placeholders, ", ")
    return placeholderString, nil
}

// Helper function to join strings with a separator
func join(elements []string, sep string) string {
    if len(elements) == 0 {
        return ""
    }
    result := elements[0]
    for _, e := range elements[1:] {
        result += sep + e
    }
    return result
}

// CheckNonZeroFields takes a struct and prints each field's name and whether it has a non-zero value.
func CheckNonZeroFields(s interface{}) bool {
    v := reflect.ValueOf(s)
    t := reflect.TypeOf(s)

    // Ensure we're working with a struct
    if v.Kind() != reflect.Struct {
        fmt.Println("Expected a struct")
        return false
    }

    // Iterate over struct fields
    for i := 0; i < v.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)

        // Check if the field is exported and has a value
        if field.PkgPath != "" {
            continue // unexported field
        }

        isZero := isZeroValue(value)
        if isZero {
            return false
        }
    }
    return true

}

// isZeroValue checks if a reflect.Value is set to its zero value.
func isZeroValue(v reflect.Value) bool {
    return v.IsZero()
}

// CompareTimeFields compares two time.Time fields from a struct, allowing for a small tolerance in seconds.
func CompareTimeFields(time1, time2 time.Time, toleranceSeconds int)bool {	 
	diff := time1.Sub(time2).Seconds()
  return diff <= float64(toleranceSeconds) && diff >= -float64(toleranceSeconds)
}

// CompareStructFields compares the fields of two structs and returns true if all fields match.// It also prints out any differences found between the two structs.
func CompareStructFields[T any](a, b T)bool {
    v1 := reflect.ValueOf(a)
    v2 := reflect.ValueOf(b)

    // Ensure that both are of the same type and kind
    if v1.Type() != v2.Type() || v1.Kind() != reflect.Struct {
        fmt.Println("Type or kind mismatch.")
        return false
    }
    
    match := true
    for i := 0; i < v1.NumField(); i++ {
        fieldA := v1.Field(i)
        fieldB := v2.Field(i)
        if fieldA.Type() == reflect.TypeOf(time.Time{}) {
            match = CompareTimeFields(fieldA.Interface().(time.Time), fieldB.Interface().(time.Time), 1)
        } else { 
            if !reflect.DeepEqual(fieldA.Interface(), fieldB.Interface()) {
              fieldName := v1.Type().Field(i).Name
              fmt.Printf("Field %s does not match: %v != %v\n", fieldName, fieldA.Interface(), fieldB.Interface())
              match = false
            }
          }
        }
    return match

    }


func GetNullable(model interface{}) []string {
    var nullableFields []string
    v := reflect.ValueOf(model)
    
    if v.Kind() == reflect.Slice {
        for i := 0; i < v.Len(); i++ {
            element := v.Index(i)
            if element.Kind() == reflect.Ptr {
                element = element.Elem()
            }
            t := element.Type()
            
            for j := 0; j < element.NumField(); j++ {
                field := element.Field(j)
                fieldType := t.Field(j)
                if field.Kind() == reflect.Ptr {
                    jsonTag := fieldType.Tag.Get("json")
                    if jsonTag != "" {
                        nullableFields = append(nullableFields, jsonTag)
                    }
                }
            }
        }
    } else if v.Kind() == reflect.Struct {
        t := v.Type()
        for i := 0; i < v.NumField(); i++ {
            field := v.Field(i)
            fieldType := t.Field(i)
            if field.Kind() == reflect.Ptr {
                jsonTag := fieldType.Tag.Get("json")
                if jsonTag != "" {
                    nullableFields = append(nullableFields, jsonTag)
                }
            }
        }
    }
    return nullableFields
}

func GetNullablePatch(model interface{}) []string {
    var nullableFields []string
    v := reflect.ValueOf(model)
    
    if v.Kind() == reflect.Slice {
        for i := 0; i < v.Len(); i++ {
            element := v.Index(i)
            if element.Kind() == reflect.Ptr {
                element = element.Elem()
            }
            t := element.Type()
            
            for j := 0; j < element.NumField(); j++ {
                fieldType := t.Field(j)
                jsonTag := fieldType.Tag.Get("json")
                requiredTag := fieldType.Tag.Get("required")
                if jsonTag != "" && requiredTag != "true"{
                    nullableFields = append(nullableFields, jsonTag)
                }
            }
        }
    } else if v.Kind() == reflect.Struct {
        t := v.Type()
        for i := 0; i < v.NumField(); i++ {
            fieldType := t.Field(i)
            jsonTag := fieldType.Tag.Get("json")
            requiredTag := fieldType.Tag.Get("required")
            if jsonTag != "" && requiredTag != "true"{
                nullableFields = append(nullableFields, jsonTag)
            }
        }
    }
    return nullableFields
}


func GetRequired(model interface{}) []string {
    var requiredFields []string
    v := reflect.ValueOf(model)
    
    if v.Kind() == reflect.Slice {
        for i := 0; i < v.Len(); i++ {
            element := v.Index(i)
            if element.Kind() == reflect.Ptr {
                element = element.Elem()
            }
            t := element.Type()
            
            for j := 0; j < element.NumField(); j++ {
                field := element.Field(j)
                fieldType := t.Field(j)
                if field.Kind() != reflect.Ptr {
                    jsonTag := fieldType.Tag.Get("json")
                    if jsonTag != "" {
                        requiredFields = append(requiredFields, jsonTag)
                    }
                }
            }
        }
    } else if v.Kind() == reflect.Struct {
        t := v.Type()
        for i := 0; i < v.NumField(); i++ {
            field := v.Field(i)
            fieldType := t.Field(i)
            if field.Kind() != reflect.Ptr {
                jsonTag := fieldType.Tag.Get("json")
                if jsonTag != "" {
                    requiredFields = append(requiredFields, jsonTag)
                }
            }
        }
    }
    return requiredFields
}

func GetRequiredPatch(model interface{}) []string {
    var requiredFields []string
    v := reflect.ValueOf(model)
    
    if v.Kind() == reflect.Slice {
        for i := 0; i < v.Len(); i++ {
            element := v.Index(i)
            if element.Kind() == reflect.Ptr {
                element = element.Elem()
            }
            t := element.Type()
            
            for j := 0; j < element.NumField(); j++ {
                fieldType := t.Field(j)
                jsonTag := fieldType.Tag.Get("json")
                requiredTag := fieldType.Tag.Get("required")
                if jsonTag != "" && requiredTag == "true"{
                    requiredFields = append(requiredFields, jsonTag)
                }
            }
        }
    } else if v.Kind() == reflect.Struct {
        t := v.Type()
        for i := 0; i < v.NumField(); i++ {
            fieldType := t.Field(i)
            jsonTag := fieldType.Tag.Get("json")
            requiredTag := fieldType.Tag.Get("required")
            if jsonTag != "" && requiredTag == "true"{
                requiredFields = append(requiredFields, jsonTag)
            }
        }
    }
    return requiredFields
}


