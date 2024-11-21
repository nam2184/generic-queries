package queries

import (
	"fmt"
	"reflect"
	"testing"

	util "github.com/nam2184/generic-queries/utils"
)

func TestGenerateConstraintFromMap(t *testing.T) {
	tests := []struct {
		name           string
		data           map[string]interface{}
		expectedString string
		expectedValues []interface{}
	}{
		{
			name: "Multiple fields",
			data: map[string]interface{}{
				"ID":  123,
				"org": "example_org",
			},
			expectedString: "ID = $1 AND org = $2",
			expectedValues: []interface{}{123, "example_org"},
		},
		{
			name: "Single field",
			data: map[string]interface{}{
				"status": "active",
			},
			expectedString: "status = $1",
			expectedValues: []interface{}{"active"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := GenerateConstraintFromMap(tt.data)

			// Check if the constraint string matches the expected output
			if constraint.constraint != tt.expectedString {
				t.Errorf("expected constraint string: %q, got: %q", tt.expectedString, constraint.constraint)
			}

			// Check if the values match the expected values
			if !reflect.DeepEqual(constraint.values, tt.expectedValues) {
				t.Errorf("expected values: %v, got: %v", tt.expectedValues, constraint.values)
			}
		})
	}
}

func TestQueryGeneration(t *testing.T) {
	// Sample input data
	tableName := "example_table"
	args := map[string]interface{}{
		"status": "active",
		"type":   "admin",
	}
	limit := 10
	skip := 20
	sortBy := "created_at"
	order := "ASC"


	// Sample constraint
	constraint := Constraint{
		constraint: "ID = $1 AND org = $2",
		values:     []interface{}{123, "example_org"},
	}

	number, _ := constraint.GetFinalPlaceholder() 
	filters, filterValues := util.GenerateFilterString(args, number, sortBy, order, limit, skip) // Adjust starting index
	// Constructing the query
	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE %s AND %s ORDER BY $%d %s LIMIT $%d OFFSET $%d",
		tableName,
		constraint.constraint,
		filters,
		number+1,
		order,
		number+2,
		number+3,
	)

	// Combine constraint values, sorting, limit, skip, and then filter values
  fullArgs := append(constraint.values, filterValues...)
	// Expected results
	expectedQuery := fmt.Sprintf(
		"SELECT * FROM %s WHERE %s AND %s ORDER BY $%d %s LIMIT $%d OFFSET $%d",
		tableName,
		"ID = $1 AND org = $2",
		"status = $6 AND type = $7",
		number+1,
		order,
		number+2,
		number+3,
	)

	expectedArgs := []interface{}{123, "example_org", sortBy, order, limit, skip, "active", "admin"}

	// Verify query format
	if query != expectedQuery {
		t.Errorf("Expected query: %q, got: %q", expectedQuery, query)
	}

	// Verify arguments
	for i, arg := range fullArgs {
		if arg != expectedArgs[i] {
			t.Errorf("At arg %d, expected: %v, got: %v", i+1, expectedArgs[i], arg)
		}
	}
}

