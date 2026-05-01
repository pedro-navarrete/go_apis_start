package product

import "context"

// Repository define la interfaz para el repositorio de productos
type Repository interface {
	// Create crea un nuevo producto
	Create(ctx context.Context, product *Product) error

	// GetByID obtiene un producto por su ID
	GetByID(ctx context.Context, id string) (*Product, error)

	// GetAll obtiene todos los productos con paginación
	GetAll(ctx context.Context, limit, offset int) ([]Product, int64, error)

	// Update actualiza un producto existente
	Update(ctx context.Context, product *Product) error

	// Delete elimina un producto por su ID
	Delete(ctx context.Context, id string) error
}
