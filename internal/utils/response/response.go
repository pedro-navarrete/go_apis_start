// Package response proporciona respuestas HTTP estandarizadas
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response estructura base para todas las respuestas
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginatedResponse respuesta con paginación
type PaginatedResponse struct {
	Success bool         `json:"success"`
	Data    interface{}  `json:"data"`
	Meta    PaginateMeta `json:"meta"`
}

// PaginateMeta metadatos de paginación
type PaginateMeta struct {
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Page   int   `json:"page"`
}

// Success envía una respuesta exitosa
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error envía una respuesta de error
func Error(c *gin.Context, statusCode int, message string, err error) {
	resp := Response{
		Success: false,
		Message: message,
	}
	if err != nil {
		resp.Error = err.Error()
	}
	c.JSON(statusCode, resp)
}

// Paginated envía una respuesta paginada
func Paginated(c *gin.Context, data interface{}, total int64, limit, offset int) {
	page := 1
	if limit > 0 {
		page = (offset/limit) + 1
	}
	c.JSON(http.StatusOK, PaginatedResponse{
		Success: true,
		Data:    data,
		Meta: PaginateMeta{
			Total:  total,
			Limit:  limit,
			Offset: offset,
			Page:   page,
		},
	})
}

// Created envía respuesta 201 Created
func Created(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusCreated, message, data)
}

// NoContent envía respuesta 204 No Content
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// BadRequest envía respuesta 400 Bad Request
func BadRequest(c *gin.Context, message string, err error) {
	Error(c, http.StatusBadRequest, message, err)
}

// Unauthorized envía respuesta 401 Unauthorized
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message, nil)
}

// NotFound envía respuesta 404 Not Found
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, nil)
}

// Conflict envía respuesta 409 Conflict
func Conflict(c *gin.Context, message string, err error) {
	Error(c, http.StatusConflict, message, err)
}

// InternalServerError envía respuesta 500 Internal Server Error
func InternalServerError(c *gin.Context, message string, err error) {
	Error(c, http.StatusInternalServerError, message, err)
}
