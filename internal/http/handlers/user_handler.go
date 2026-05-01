package handlers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/pedro-navarrete/go_apis_start/internal/domain/user"
	apperrors "github.com/pedro-navarrete/go_apis_start/pkg/errors"
	"github.com/pedro-navarrete/go_apis_start/internal/utils/response"
	"github.com/pedro-navarrete/go_apis_start/internal/utils/validator"
)

// UserHandler maneja los endpoints de usuarios
type UserHandler struct {
	service   user.Service
	jwtSecret string
}

// NewUserHandler crea una nueva instancia del handler de usuarios
func NewUserHandler(service user.Service, jwtSecret string) *UserHandler {
	return &UserHandler{
		service:   service,
		jwtSecret: jwtSecret,
	}
}

// Create crea un nuevo usuario
// POST /api/users
func (h *UserHandler) Create(c *gin.Context) {
	var req user.CreateUserRequest

	// Parsear el body JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "datos inválidos", err)
		return
	}

	// Validar los datos
	if err := validator.Validate(req); err != nil {
		response.BadRequest(c, "error de validación", err)
		return
	}

	// Crear el usuario
	newUser, err := h.service.CreateUser(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, apperrors.ErrConflict) {
			response.Conflict(c, "el usuario ya existe", err)
			return
		}
		response.InternalServerError(c, "error al crear usuario", err)
		return
	}

	response.Created(c, "usuario creado exitosamente", newUser)
}

// Login autentica un usuario y retorna un token JWT
// POST /api/auth/login
func (h *UserHandler) Login(c *gin.Context) {
	var req user.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "datos inválidos", err)
		return
	}

	if err := validator.Validate(req); err != nil {
		response.BadRequest(c, "error de validación", err)
		return
	}

	loginResp, err := h.service.Login(c.Request.Context(), req, h.jwtSecret)
	if err != nil {
		if errors.Is(err, apperrors.ErrUnauthorized) {
			response.Unauthorized(c, "usuario o contraseña incorrectos")
			return
		}
		response.InternalServerError(c, "error al iniciar sesión", err)
		return
	}

	response.Success(c, 200, "login exitoso", loginResp)
}

// GetByID obtiene un usuario por su ID
// GET /api/users/:id
func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	u, err := h.service.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			response.NotFound(c, "usuario no encontrado")
			return
		}
		response.InternalServerError(c, "error al obtener usuario", err)
		return
	}

	response.Success(c, 200, "usuario obtenido", u)
}

// List lista todos los usuarios con paginación
// GET /api/users
func (h *UserHandler) List(c *gin.Context) {
	// Parámetros de paginación
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	users, total, err := h.service.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(c, "error al listar usuarios", err)
		return
	}

	response.Paginated(c, users, total, limit, offset)
}

// Update actualiza un usuario
// PUT /api/users/:id
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req user.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "datos inválidos", err)
		return
	}

	if err := validator.Validate(req); err != nil {
		response.BadRequest(c, "error de validación", err)
		return
	}

	updatedUser, err := h.service.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			response.NotFound(c, "usuario no encontrado")
			return
		}
		response.InternalServerError(c, "error al actualizar usuario", err)
		return
	}

	response.Success(c, 200, "usuario actualizado exitosamente", updatedUser)
}

// Delete elimina un usuario
// DELETE /api/users/:id
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteUser(c.Request.Context(), id); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			response.NotFound(c, "usuario no encontrado")
			return
		}
		response.InternalServerError(c, "error al eliminar usuario", err)
		return
	}

	response.NoContent(c)
}
