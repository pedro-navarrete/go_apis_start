package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/pedro-navarrete/go_apis_start/internal/domain/product"
)

// productSQLServerRepository implementación del repositorio de productos para SQL Server
type productSQLServerRepository struct {
	db *gorm.DB
}

// NewProductSQLServerRepository crea una nueva instancia del repositorio de productos para SQL Server
func NewProductSQLServerRepository(db *gorm.DB) product.Repository {
	return &productSQLServerRepository{db: db}
}

// Create crea un nuevo producto en SQL Server
func (r *productSQLServerRepository) Create(ctx context.Context, p *product.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

// GetByID obtiene un producto por su ID de SQL Server
func (r *productSQLServerRepository) GetByID(ctx context.Context, id string) (*product.Product, error) {
	var p product.Product
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&p)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return &p, result.Error
}

// GetAll obtiene todos los productos con paginación de SQL Server
func (r *productSQLServerRepository) GetAll(ctx context.Context, limit, offset int) ([]product.Product, int64, error) {
	var products []product.Product
	var total int64

	if err := r.db.WithContext(ctx).Model(&product.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	result := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&products)
	return products, total, result.Error
}

// Update actualiza un producto en SQL Server
func (r *productSQLServerRepository) Update(ctx context.Context, p *product.Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}

// Delete elimina un producto de SQL Server
func (r *productSQLServerRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&product.Product{}).Error
}
