// Package user define el dominio del módulo de usuarios
package user

import "time"

// User representa el modelo de usuario en la base de datos
type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null" validate:"required,min=3,max=50"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Password  string    `json:"-" gorm:"not null"` // omitido en JSON por seguridad
	FullName  string    `json:"full_name" gorm:"not null" validate:"required"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName define el nombre de la tabla en SQL Server
func (User) TableName() string {
	return "users"
}

// CreateUserRequest datos requeridos para crear un usuario
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
}

// UpdateUserRequest datos para actualizar un usuario (todos opcionales)
type UpdateUserRequest struct {
	Email    string `json:"email" validate:"omitempty,email"`
	FullName string `json:"full_name" validate:"omitempty"`
}

// LoginRequest datos para el login
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse respuesta después de un login exitoso
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
