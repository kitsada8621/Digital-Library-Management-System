package services

import (
	"dlms/dtos"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type IValidationService interface {
	Validate(err error) []dtos.CustomError
}

type ValidationServiceImpl struct{}

func ValidationService() IValidationService {
	return &ValidationServiceImpl{}
}

func (s *ValidationServiceImpl) Validate(err error) []dtos.CustomError {
	results := []dtos.CustomError{}
	if validateErrs, ok := err.(validator.ValidationErrors); ok {
		for _, row := range validateErrs {
			results = append(results, dtos.CustomError{
				Key:   strings.ToLower(row.Field()),
				Value: fmt.Sprintf("The %s is %s", strings.ToLower(row.Field()), row.Tag()),
			})
		}
	}

	return results
}
