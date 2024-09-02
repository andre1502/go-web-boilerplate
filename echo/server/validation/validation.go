package validation

import (
	"boilerplate/utils"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Validation struct {
	validator *validator.Validate
}

func NewValidation() *Validation {
	validation := &Validation{
		validator: validator.New(),
	}

	validation.validator.RegisterTagNameFunc(validation.GetJsonTagName())
	validation.validator.RegisterValidation("empty_string", validation.EmptyString())

	return validation
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

func (v *Validation) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

func (v *Validation) ValidateRequest(c echo.Context, i interface{}) error {
	if err := c.Bind(i); err != nil {
		return err
	}

	if err := c.Validate(i); err != nil {
		return err
	}

	return nil
}
