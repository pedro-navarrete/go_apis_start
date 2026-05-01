// Package database maneja las conexiones a las bases de datos
package database

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/pedro-navarrete/go_apis_start/internal/config"
	applogger "github.com/pedro-navarrete/go_apis_start/internal/utils/logger"
)

// SQLServerDB wraps the GORM DB connection for SQL Server
type SQLServerDB struct {
	DB *gorm.DB
}

// NewSQLServerConnection crea una nueva conexión a SQL Server usando GORM
func NewSQLServerConnection(cfg *config.SQLServerConfig) (*SQLServerDB, error) {
	// Construir el connection string de SQL Server
	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	// Configurar el nivel de log de GORM según entorno
	gormLogger := logger.Default.LogMode(logger.Silent)

	// Conectar a la base de datos
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("error conectando a SQL Server: %w", err)
	}

	// Configurar el connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo instancia SQL: %w", err)
	}

	// Configuración del pool de conexiones
	sqlDB.SetMaxOpenConns(25)                 // Máximo de conexiones abiertas
	sqlDB.SetMaxIdleConns(10)                 // Máximo de conexiones inactivas
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Tiempo máximo de vida de una conexión

	applogger.Info("✅ Conexión a SQL Server establecida")

	return &SQLServerDB{DB: db}, nil
}

// Ping verifica que la conexión a SQL Server esté activa
func (s *SQLServerDB) Ping() error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Close cierra la conexión a SQL Server
func (s *SQLServerDB) Close() error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// AutoMigrate ejecuta las migraciones automáticas de GORM
func (s *SQLServerDB) AutoMigrate(models ...interface{}) error {
	return s.DB.AutoMigrate(models...)
}
