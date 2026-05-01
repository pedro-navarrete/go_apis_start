package user

import "context"

// Repository define la interfaz para el repositorio de usuarios
// Permite intercambiar implementaciones (SQL Server, MongoDB, etc.)
type Repository interface {
	// Create crea un nuevo usuario
	Create(ctx context.Context, user *User) error

	// GetByID obtiene un usuario por su ID
	GetByID(ctx context.Context, id string) (*User, error)

	// GetByUsername obtiene un usuario por su nombre de usuario
	GetByUsername(ctx context.Context, username string) (*User, error)

	// GetByEmail obtiene un usuario por su email
	GetByEmail(ctx context.Context, email string) (*User, error)

	// GetAll obtiene todos los usuarios con paginación
	GetAll(ctx context.Context, limit, offset int) ([]User, int64, error)

	// Update actualiza un usuario existente
	Update(ctx context.Context, user *User) error

	// Delete elimina un usuario por su ID
	Delete(ctx context.Context, id string) error
}
