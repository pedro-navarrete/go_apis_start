package handlers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/pedro-navarrete/go_apis_start/internal/domain/product"
	apperrors "github.com/pedro-navarrete/go_apis_start/pkg/errors"
	"github.com/pedro-navarrete/go_apis_start/internal/utils/response"
	"github.com/pedro-navarrete/go_apis_start/internal/utils/validator"
)

// ProductHandler maneja los endpoints de productos
type ProductHandler struct {
	service product.Service
}

// NewProductHandler crea una nueva instancia del handler de productos
func NewProductHandler(service product.Service) *ProductHandler {
	return &ProductHandler{service: service}
}

// Create crea un nuevo producto
// POST /api/products
func (h *ProductHandler) Create(c *gin.Context) {
	var req product.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "datos inválidos", err)
		return
	}

	if err := validator.Validate(req); err != nil {
		response.BadRequest(c, "error de validación", err)
		return
	}

	newProduct, err := h.service.CreateProduct(c.Request.Context(), req)
	if err != nil {
		response.InternalServerError(c, "error al crear producto", err)
		return
	}

	response.Created(c, "producto creado exitosamente", newProduct)
}

// GetByID obtiene un producto por su ID
// GET /api/products/:id
func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	p, err := h.service.GetProductByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			response.NotFound(c, "producto no encontrado")
			return
		}
		response.InternalServerError(c, "error al obtener producto", err)
		return
	}

	response.Success(c, 200, "producto obtenido", p)
}

// List lista todos los productos con paginación
// GET /api/products
func (h *ProductHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	products, total, err := h.service.ListProducts(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(c, "error al listar productos", err)
		return
	}

	response.Paginated(c, products, total, limit, offset)
}

// Update actualiza un producto
// PUT /api/products/:id
func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req product.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "datos inválidos", err)
		return
	}

	if err := validator.Validate(req); err != nil {
		response.BadRequest(c, "error de validación", err)
		return
	}

	updatedProduct, err := h.service.UpdateProduct(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			response.NotFound(c, "producto no encontrado")
			return
		}
		response.InternalServerError(c, "error al actualizar producto", err)
		return
	}

	response.Success(c, 200, "producto actualizado exitosamente", updatedProduct)
}

// Delete elimina un producto
// DELETE /api/products/:id
func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteProduct(c.Request.Context(), id); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			response.NotFound(c, "producto no encontrado")
			return
		}
		response.InternalServerError(c, "error al eliminar producto", err)
		return
	}

	response.NoContent(c)
}
