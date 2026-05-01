// Package routes registra todas las rutas de la aplicación
package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/pedro-navarrete/go_apis_start/internal/http/handlers"
	"github.com/pedro-navarrete/go_apis_start/internal/http/middleware"
)

// Dependencies contiene todas las dependencias necesarias para las rutas
type Dependencies struct {
	UserHandler    *handlers.UserHandler
	ProductHandler *handlers.ProductHandler
	JWTSecret      string
	BasicAuthUsers map[string]string
}

// SetupRoutes configura todas las rutas de la API
func SetupRoutes(router *gin.Engine, deps *Dependencies) {
	// Grupo base de la API
	api := router.Group("/api")

	// Health check
	api.GET("/health", handlers.HealthHandler)

	// Autenticación
	auth := api.Group("/auth")
	{
		auth.POST("/login", deps.UserHandler.Login)
	}

	// Usuarios
	users := api.Group("/users")
	{
		// Ruta pública: crear usuario
		users.POST("", deps.UserHandler.Create)

		// Rutas protegidas con JWT
		protected := users.Group("")
		protected.Use(middleware.JWTAuth(deps.JWTSecret))
		{
			protected.GET("", deps.UserHandler.List)
			protected.GET("/:id", deps.UserHandler.GetByID)
			protected.PUT("/:id", deps.UserHandler.Update)
			protected.DELETE("/:id", deps.UserHandler.Delete)
		}
	}

	// Productos
	products := api.Group("/products")
	{
		// Rutas públicas
		products.GET("", deps.ProductHandler.List)
		products.GET("/:id", deps.ProductHandler.GetByID)

		// Crear producto: requiere Basic Auth
		products.POST("", middleware.BasicAuth(deps.BasicAuthUsers), deps.ProductHandler.Create)

		// Rutas protegidas con JWT
		protected := products.Group("")
		protected.Use(middleware.JWTAuth(deps.JWTSecret))
		{
			protected.PUT("/:id", deps.ProductHandler.Update)
			protected.DELETE("/:id", deps.ProductHandler.Delete)
		}
	}
}
