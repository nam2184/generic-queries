package util

import (
	"fmt"
	"reflect"
	"time"
)

// IsZeroField checks if a given field in a struct is set to its zero value.
func IsZeroField(v interface{}, fieldName string) bool {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
    fmt.Println("Not a struct")
		return false
	}

	field := val.FieldByName(fieldName)
	if !field.IsValid() {
    fmt.Printf("No such field: %s in struct", fieldName)
		return false	
  }

	return reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface())
}

func GetZero[T any]() T {
    var result T
    return result
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
            if !isZeroValue(value){
              placeholders = append(placeholders, fmt.Sprintf(":%s", dbTag))
              params = append(params, dbTag)
        }
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
             if !isZeroValue(value){
              placeholders = append(placeholders, fmt.Sprintf(":%s", dbTag))
              params = append(params, dbTag)
            }
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
            placeholders = append(placeholders, fmt.Sprintf(":%s:%s", dbTag, dbTag))
            params = append(params, dbTag)
        default:
            if !isZeroValue(value){
              placeholders = append(placeholders, fmt.Sprintf(":%s:%s", dbTag, dbTag))
              params = append(params, dbTag)
        }
        }
    }

    // Construct the SQL parameter placeholder string
    placeholderString := fmt.Sprintf("%s", join(params, ", "))
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
	// Calculate the difference in seconds
	diff := time1.Sub(time2).Seconds()
	// Compare with tolerance
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