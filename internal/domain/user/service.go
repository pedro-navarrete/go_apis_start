package user

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	apperrors "github.com/pedro-navarrete/go_apis_start/pkg/errors"
	"github.com/pedro-navarrete/go_apis_start/internal/utils/password"
)

// Service define la interfaz del servicio de usuarios
type Service interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)
	Login(ctx context.Context, req LoginRequest, jwtSecret string) (*LoginResponse, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, id string, req UpdateUserRequest) (*User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, limit, offset int) ([]User, int64, error)
}

// service implementación del servicio de usuarios
type service struct {
	repo          Repository
	jwtExpiration time.Duration
}

// NewService crea una nueva instancia del servicio de usuarios
func NewService(repo Repository, jwtExpiration time.Duration) Service {
	return &service{
		repo:          repo,
		jwtExpiration: jwtExpiration,
	}
}

// CreateUser crea un nuevo usuario con contraseña hasheada
func (s *service) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
	// Verificar si el username ya existe
	existing, _ := s.repo.GetByUsername(ctx, req.Username)
	if existing != nil {
		return nil, apperrors.ErrConflict
	}

	// Verificar si el email ya existe
	existing, _ = s.repo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, apperrors.ErrConflict
	}

	// Hashear la contraseña
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("error hasheando contraseña: %w", err)
	}

	// Crear el usuario
	user := &User{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		FullName:  req.FullName,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("error creando usuario: %w", err)
	}

	return user, nil
}

// Login valida las credenciales y genera un token JWT
func (s *service) Login(ctx context.Context, req LoginRequest, jwtSecret string) (*LoginResponse, error) {
	// Buscar usuario por username
	user, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, apperrors.ErrUnauthorized
	}

	// Verificar la contraseña
	if !password.Verify(user.Password, req.Password) {
		return nil, apperrors.ErrUnauthorized
	}

	// Generar token JWT
	token, err := generateJWT(user, jwtSecret, s.jwtExpiration)
	if err != nil {
		return nil, fmt.Errorf("error generando token: %w", err)
	}

	return &LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

// GetUserByID obtiene un usuario por su ID
func (s *service) GetUserByID(ctx context.Context, id string) (*User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	return user, nil
}

// UpdateUser actualiza los datos de un usuario
func (s *service) UpdateUser(ctx context.Context, id string, req UpdateUserRequest) (*User, error) {
	// Obtener usuario existente
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}

	// Actualizar campos si se proporcionaron
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("error actualizando usuario: %w", err)
	}

	return user, nil
}

// DeleteUser elimina un usuario por su ID
func (s *service) DeleteUser(ctx context.Context, id string) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apperrors.ErrNotFound
	}

	return s.repo.Delete(ctx, id)
}

// ListUsers obtiene la lista de usuarios con paginación
func (s *service) ListUsers(ctx context.Context, limit, offset int) ([]User, int64, error) {
	return s.repo.GetAll(ctx, limit, offset)
}

// JWTClaims define los claims del token JWT
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// generateJWT genera un token JWT para el usuario
func generateJWT(user *User, secret string, expiration time.Duration) (string, error) {
	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
