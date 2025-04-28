package validator

import (
	"github.com/IlhamSetiaji/julong-notification-be/config"
	"github.com/go-playground/validator/v10"
)

type validatorV10 struct {
	ValidatorV10 *validator.Validate
}

func NewValidatorV10(conf *config.Config) Validator {
	validate := validator.New()
	validate.RegisterValidation("application", func(fl validator.FieldLevel) bool {
		templateType := fl.Field().String()
		switch templateType {
		case "MANPOWER", "RECRUITMENT", "ONBOARDING":
			return true
		default:
			return false
		}
	})
	return &validatorV10{
		ValidatorV10: validate,
	}
}

func (v *validatorV10) GetValidator() *validator.Validate {
	return v.ValidatorV10
}
