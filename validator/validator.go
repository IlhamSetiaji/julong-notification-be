package validator

import "github.com/go-playground/validator/v10"

type Validator interface {
	GetValidator() *validator.Validate
}
