// Package validator proporciona validación de structs con mensajes de error descriptivos
package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator instancia del validador
var validate *validator.Validate

// init inicializa el validador
func init() {
	validate = validator.New()
}

// Validate valida un struct y retorna errores formateados
func Validate(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		return formatValidationErrors(err)
	}
	return nil
}

// formatValidationErrors formatea los errores de validación en mensajes legibles
func formatValidationErrors(err error) error {
	var messages []string

	for _, e := range err.(validator.ValidationErrors) {
		messages = append(messages, formatError(e))
	}

	return fmt.Errorf("%s", strings.Join(messages, "; "))
}

// formatError formatea un error de validación individual
func formatError(e validator.FieldError) string {
	field := strings.ToLower(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("el campo '%s' es requerido", field)
	case "email":
		return fmt.Sprintf("el campo '%s' debe ser un email válido", field)
	case "min":
		return fmt.Sprintf("el campo '%s' debe tener mínimo %s caracteres", field, e.Param())
	case "max":
		return fmt.Sprintf("el campo '%s' debe tener máximo %s caracteres", field, e.Param())
	case "gt":
		return fmt.Sprintf("el campo '%s' debe ser mayor a %s", field, e.Param())
	case "gte":
		return fmt.Sprintf("el campo '%s' debe ser mayor o igual a %s", field, e.Param())
	case "lt":
		return fmt.Sprintf("el campo '%s' debe ser menor a %s", field, e.Param())
	case "lte":
		return fmt.Sprintf("el campo '%s' debe ser menor o igual a %s", field, e.Param())
	case "oneof":
		return fmt.Sprintf("el campo '%s' debe ser uno de: %s", field, e.Param())
	default:
		return fmt.Sprintf("el campo '%s' es inválido (%s)", field, e.Tag())
	}
}
