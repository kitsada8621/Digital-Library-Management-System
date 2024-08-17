package utils

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

func CustomEqField(fl validator.FieldLevel) bool {

	field := fl.Field()
	otherFieldName := fl.Param()

	parent := fl.Top()
	otherField := reflect.Indirect(parent).FieldByName(otherFieldName)

	if !otherField.IsValid() {
		return false
	}

	return field.Interface() == otherField.Interface()
}
