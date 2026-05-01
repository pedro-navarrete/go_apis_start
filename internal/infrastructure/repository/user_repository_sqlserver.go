// Package repository contiene las implementaciones concretas de los repositorios
package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/pedro-navarrete/go_apis_start/internal/domain/user"
)

// userSQLServerRepository implementación del repositorio de usuarios para SQL Server
type userSQLServerRepository struct {
	db *gorm.DB
}

// NewUserSQLServerRepository crea una nueva instancia del repositorio de usuarios para SQL Server
func NewUserSQLServerRepository(db *gorm.DB) user.Repository {
	return &userSQLServerRepository{db: db}
}

// Create crea un nuevo usuario en SQL Server
func (r *userSQLServerRepository) Create(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

// GetByID obtiene un usuario por su ID de SQL Server
func (r *userSQLServerRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	var u user.User
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&u)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return &u, result.Error
}

// GetByUsername obtiene un usuario por su username de SQL Server
func (r *userSQLServerRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	var u user.User
	result := r.db.WithContext(ctx).Where("username = ?", username).First(&u)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return &u, result.Error
}

// GetByEmail obtiene un usuario por su email de SQL Server
func (r *userSQLServerRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&u)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	return &u, result.Error
}

// GetAll obtiene todos los usuarios con paginación de SQL Server
func (r *userSQLServerRepository) GetAll(ctx context.Context, limit, offset int) ([]user.User, int64, error) {
	var users []user.User
	var total int64

	// Contar total de registros
	if err := r.db.WithContext(ctx).Model(&user.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Obtener registros con paginación
	result := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users)
	return users, total, result.Error
}

// Update actualiza un usuario en SQL Server
func (r *userSQLServerRepository) Update(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

// Delete elimina un usuario de SQL Server
func (r *userSQLServerRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&user.User{}).Error
}
