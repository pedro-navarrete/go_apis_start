// Package main es el punto de entrada de la aplicación
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/pedro-navarrete/go_apis_start/internal/config"
	"github.com/pedro-navarrete/go_apis_start/internal/domain/product"
	"github.com/pedro-navarrete/go_apis_start/internal/domain/user"
	"github.com/pedro-navarrete/go_apis_start/internal/http/handlers"
	"github.com/pedro-navarrete/go_apis_start/internal/http/middleware"
	"github.com/pedro-navarrete/go_apis_start/internal/http/routes"
	"github.com/pedro-navarrete/go_apis_start/internal/infrastructure/database"
	"github.com/pedro-navarrete/go_apis_start/internal/infrastructure/repository"
	applogger "github.com/pedro-navarrete/go_apis_start/internal/utils/logger"
)

func main() {
	// 1. Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error cargando configuración: %v\n", err)
		os.Exit(1)
	}

	// 2. Inicializar logger
	if err := applogger.Init(cfg.Log.Level, cfg.Server.Environment); err != nil {
		fmt.Printf("Error inicializando logger: %v\n", err)
		os.Exit(1)
	}
	defer applogger.Sync()

	applogger.Info("🚀 Iniciando go_apis_start",
		zap.String("environment", cfg.Server.Environment),
		zap.String("port", cfg.Server.Port),
	)

	// 3. Conectar a bases de datos
	var userRepo user.Repository
	var productRepo product.Repository

	// Intentar conectar a SQL Server con retry
	var sqlDB *database.SQLServerDB
	for i := 0; i < 5; i++ {
		sqlDB, err = database.NewSQLServerConnection(&cfg.SQLServer)
		if err == nil {
			break
		}
		applogger.Warn("Reintentando conexión a SQL Server...",
			zap.Int("intento", i+1),
			zap.Error(err),
		)
		time.Sleep(3 * time.Second)
	}

	// Conectar a MongoDB
	var mongoDB *database.MongoDB
	for i := 0; i < 5; i++ {
		mongoDB, err = database.NewMongoDBConnection(&cfg.MongoDB)
		if err == nil {
			break
		}
		applogger.Warn("Reintentando conexión a MongoDB...",
			zap.Int("intento", i+1),
			zap.Error(err),
		)
		time.Sleep(3 * time.Second)
	}

	// 4. Configurar repositorios según el tipo de BD configurado
	if cfg.DBTypes.UserDBType == "mongodb" && mongoDB != nil {
		userRepo = repository.NewUserMongoRepository(mongoDB.Collection("users"))
		applogger.Info("📦 Usuario usando MongoDB")
	} else if sqlDB != nil {
		// Auto-migrar modelos en SQL Server
		if err := sqlDB.AutoMigrate(&user.User{}); err != nil {
			applogger.Warn("Error en auto-migración de usuarios", zap.Error(err))
		}
		userRepo = repository.NewUserSQLServerRepository(sqlDB.DB)
		applogger.Info("📦 Usuario usando SQL Server")
	} else {
		applogger.Fatal("No se pudo conectar a ninguna base de datos para usuarios")
		os.Exit(1)
	}

	if cfg.DBTypes.ProductDBType == "mongodb" && mongoDB != nil {
		productRepo = repository.NewProductMongoRepository(mongoDB.Collection("products"))
		applogger.Info("📦 Producto usando MongoDB")
	} else if sqlDB != nil {
		// Auto-migrar modelos en SQL Server
		if err := sqlDB.AutoMigrate(&product.Product{}); err != nil {
			applogger.Warn("Error en auto-migración de productos", zap.Error(err))
		}
		productRepo = repository.NewProductSQLServerRepository(sqlDB.DB)
		applogger.Info("📦 Producto usando SQL Server")
	} else {
		applogger.Fatal("No se pudo conectar a ninguna base de datos para productos")
		os.Exit(1)
	}

	// 5. Crear servicios
	userService := user.NewService(userRepo, cfg.JWT.Expiration)
	productService := product.NewService(productRepo)

	// 6. Crear handlers
	userHandler := handlers.NewUserHandler(userService, cfg.JWT.Secret)
	productHandler := handlers.NewProductHandler(productService)

	// 7. Configurar el router Gin
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// 8. Registrar middlewares globales
	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// 9. Registrar rutas
	routes.SetupRoutes(router, &routes.Dependencies{
		UserHandler:    userHandler,
		ProductHandler: productHandler,
		JWTSecret:      cfg.JWT.Secret,
		BasicAuthUsers: cfg.BasicAuth.Users,
	})

	// 10. Configurar el servidor HTTP
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 11. Iniciar servidor en goroutine para graceful shutdown
	go func() {
		applogger.Info(fmt.Sprintf("🌐 Servidor escuchando en http://localhost:%s", cfg.Server.Port))
		applogger.Info("📋 Endpoints disponibles:")
		applogger.Info("  GET    /api/health")
		applogger.Info("  POST   /api/auth/login")
		applogger.Info("  POST   /api/users")
		applogger.Info("  GET    /api/users          (JWT)")
		applogger.Info("  GET    /api/users/:id      (JWT)")
		applogger.Info("  PUT    /api/users/:id      (JWT)")
		applogger.Info("  DELETE /api/users/:id      (JWT)")
		applogger.Info("  GET    /api/products")
		applogger.Info("  GET    /api/products/:id")
		applogger.Info("  POST   /api/products       (Basic Auth)")
		applogger.Info("  PUT    /api/products/:id   (JWT)")
		applogger.Info("  DELETE /api/products/:id   (JWT)")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			applogger.Fatal("Error iniciando servidor", zap.Error(err))
		}
	}()

	// 12. Esperar señal de terminación (graceful shutdown)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	applogger.Info("🛑 Iniciando shutdown graceful...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		applogger.Error("Error en shutdown del servidor", zap.Error(err))
	}

	// Cerrar conexiones a bases de datos
	if sqlDB != nil {
		if err := sqlDB.Close(); err != nil {
			applogger.Error("Error cerrando SQL Server", zap.Error(err))
		}
	}
	if mongoDB != nil {
		if err := mongoDB.Close(ctx); err != nil {
			applogger.Error("Error cerrando MongoDB", zap.Error(err))
		}
	}

	applogger.Info("✅ Servidor detenido correctamente")
}
