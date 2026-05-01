package product

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	apperrors "github.com/pedro-navarrete/go_apis_start/pkg/errors"
)

// Service define la interfaz del servicio de productos
type Service interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (*Product, error)
	GetProductByID(ctx context.Context, id string) (*Product, error)
	UpdateProduct(ctx context.Context, id string, req UpdateProductRequest) (*Product, error)
	DeleteProduct(ctx context.Context, id string) error
	ListProducts(ctx context.Context, limit, offset int) ([]Product, int64, error)
}

// service implementación del servicio de productos
type service struct {
	repo Repository
}

// NewService crea una nueva instancia del servicio de productos
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateProduct crea un nuevo producto
func (s *service) CreateProduct(ctx context.Context, req CreateProductRequest) (*Product, error) {
	product := &Product{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("error creando producto: %w", err)
	}

	return product, nil
}

// GetProductByID obtiene un producto por su ID
func (s *service) GetProductByID(ctx context.Context, id string) (*Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	return product, nil
}

// UpdateProduct actualiza los datos de un producto
func (s *service) UpdateProduct(ctx context.Context, id string, req UpdateProductRequest) (*Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	product.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("error actualizando producto: %w", err)
	}

	return product, nil
}

// DeleteProduct elimina un producto por su ID
func (s *service) DeleteProduct(ctx context.Context, id string) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperrors.ErrNotFound
	}
	return s.repo.Delete(ctx, id)
}

// ListProducts obtiene la lista de productos con paginación
func (s *service) ListProducts(ctx context.Context, limit, offset int) ([]Product, int64, error) {
	return s.repo.GetAll(ctx, limit, offset)
}
