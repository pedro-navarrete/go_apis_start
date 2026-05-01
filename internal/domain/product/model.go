// Package product define el dominio del módulo de productos
package product

import "time"

// Product representa el modelo de producto en la base de datos
type Product struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null" validate:"required,min=3"`
	Description string    `json:"description"`
	Price       float64   `json:"price" gorm:"not null" validate:"required,gt=0"`
	Stock       int       `json:"stock" gorm:"default:0" validate:"gte=0"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName define el nombre de la tabla en SQL Server
func (Product) TableName() string {
	return "products"
}

// CreateProductRequest datos requeridos para crear un producto
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=3"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
}

// UpdateProductRequest datos para actualizar un producto (todos opcionales)
type UpdateProductRequest struct {
	Name        string  `json:"name" validate:"omitempty,min=3"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"omitempty,gt=0"`
	Stock       int     `json:"stock" validate:"omitempty,gte=0"`
}
