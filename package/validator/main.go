package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Validator struct {
	Validate *validator.Validate
}

func (v Validator) UnmarshallJSONValidate(c echo.Context, req any) error {
	if err := c.Bind(req); err != nil {
		return err
	}
	return v.ValidateStruct(req)
}

func (v Validator) ValidateStruct(val any) error {
	return v.Validate.Struct(val)
}

func NewValidator() *Validator {
	return &Validator{
		Validate: validator.New(),
	}
}
