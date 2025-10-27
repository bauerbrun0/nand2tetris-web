package validator

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	FieldErrors    map[string]string
	NonFieldErrors []string
	Validate       *validator.Validate
}

func NewValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterValidation("no_whitespace", func(fl validator.FieldLevel) bool {
		for _, r := range fl.Field().String() {
			if unicode.IsSpace(r) {
				return false
			}
		}
		return true
	})
	return v
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckFieldBool(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func (v *Validator) CheckFieldError(err error, key, message string) {
	if err != nil {
		v.AddFieldError(key, message)
	}
}

func (v *Validator) CheckFieldTag(field any, tag, key, message string) {
	err := v.Validate.Var(field, tag)
	if err != nil {
		v.AddFieldError(key, message)
	}
}

func (v *Validator) GetFirstFieldError() string {
	for _, msg := range v.FieldErrors {
		return msg
	}
	return ""
}
