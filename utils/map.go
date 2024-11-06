package util

import (
	"fmt"
	"strconv"
	"strings"
)

func GetKeyAndValues(data map[string]interface{}) (string, []interface{}) {
    keys := make([]string, 0, len(data))
    values := make([]interface{}, 0, len(data))
    for key, value := range data {
        keys = append(keys, key)
        values = append(values,value)
    }
    return fmt.Sprintf("%s", strings.Join(keys, ", ")), values
}

func GetValues(data map[string]interface{}) []interface{} {
    values := make([]interface{}, 0, len(data))
    for _, value := range data {
        values = append(values, value)
    }
    return values
}

func GetNumberedParameters(data map[string]interface{}) string {
    placeholders := make([]string, 0, len(data))
    for i := 1; i <= len(data); i++ {
        placeholders = append(placeholders, fmt.Sprintf("$%d", i))
    }
    return strings.Join(placeholders, ", ")
}

func GenerateFilterString(params map[string]interface{}, limit int, offset int, sort_by string) (string, []interface{}) {
    var filterString string
    var values []interface{}
    count := 3 // Start numbering placeholders from 3
    values = append(values, limit, offset, sort_by)
    for key, value := range params {
        if filterString != "" {
            filterString += " AND "
        }
        filterString += fmt.Sprintf("%s = $%d", key, count)
        values = append(values, value)
        count++
    }
    return filterString, values
}

func GetIntFromInterface(value interface{}) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true 
	case string:
		if intValue, err := strconv.Atoi(v); err == nil {
			return intValue, true
		}
	default:
		return 0, false
	}
	return 0, false // Return zero if conversion is not possible
}
