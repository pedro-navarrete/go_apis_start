// Package errors define los errores personalizados de la aplicación
package errors

import "errors"

// Errores de dominio estándar de la aplicación
var (
	// ErrNotFound se retorna cuando un recurso no es encontrado
	ErrNotFound = errors.New("resource not found")

	// ErrBadRequest se retorna cuando la solicitud es inválida
	ErrBadRequest = errors.New("bad request")

	// ErrUnauthorized se retorna cuando el usuario no está autenticado
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden se retorna cuando el usuario no tiene permisos
	ErrForbidden = errors.New("forbidden")

	// ErrConflict se retorna cuando hay un conflicto (ej: recurso ya existe)
	ErrConflict = errors.New("resource already exists")

	// ErrInternal se retorna cuando ocurre un error interno del servidor
	ErrInternal = errors.New("internal server error")
)
