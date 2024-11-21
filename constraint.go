package queries

import (
	"fmt"
	"strings"

	util "github.com/nam2184/generic-queries/utils"
)

type Constraint struct {
	constraint string
	values     []interface{}
}

func NewConstraint(constraint string, values ...interface{}) (Constraint, error) {
  new_constraint, err := GetConstraintString(constraint, values); if err != nil {
    return util.GetZero[Constraint](), err
  }
  return Constraint{constraint: new_constraint, values: values}, nil
}

func (c Constraint) GetFinalPlaceholder() (int, error) {
	if len(c.values) == 0 {
		return 0, fmt.Errorf("no values provided")
	}
	return len(c.values), nil
}


// GetConstraintString replaces placeholders with numbered placeholders
func GetConstraintString(constraint string, values ...interface{}) (string, error) {
	placeholderCount := strings.Count(constraint, "$%")
	if placeholderCount != len(values) {
		return "", fmt.Errorf("mismatch between placeholders and values count")
	}

	finalConstraint := constraint
	for i := 1; i <= len(values); i++ {
		placeholder := fmt.Sprintf("$%d", i)
		finalConstraint = strings.Replace(finalConstraint, "$%", placeholder, 1)
	}
	return finalConstraint, nil
}

func GenerateConstraintFromMap(data map[string]interface{}) Constraint {
	var placeholders []string
	var values []interface{}
	counter := 1

	for key, value := range data {
		placeholders = append(placeholders, fmt.Sprintf("%s = $%d", key, counter))
		values = append(values, value)
		counter++
	}

	constraintStr := strings.Join(placeholders, " AND ")

	return Constraint{constraint: constraintStr, values: values}
}
