package validation

import (
	"boilerplate/utils"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validation struct {
	validator *validator.Validate
}

func NewValidation() *Validation {
	return &Validation{
		validator: validator.New(),
	}
}

func (v *Validation) EmptyString() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		if str, ok := fl.Field().Interface().(string); ok {
			if utils.IsEmptyString(str) {
				return false
			}
		}

		return true
	}
}

func (v *Validation) GetJsonTagName() func(fld reflect.StructField) string {
	return func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	}
}
