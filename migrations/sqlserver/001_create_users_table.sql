-- Migración 001: Crear tabla de usuarios
-- Fecha: 2024-01-01
-- Descripción: Crea la tabla principal de usuarios con índices

CREATE TABLE users (
    id          NVARCHAR(36)  NOT NULL PRIMARY KEY,
    username    NVARCHAR(50)  NOT NULL UNIQUE,
    email       NVARCHAR(255) NOT NULL UNIQUE,
    password    NVARCHAR(255) NOT NULL,
    full_name   NVARCHAR(255) NOT NULL,
    is_active   BIT           NOT NULL DEFAULT 1,
    created_at  DATETIME2     NOT NULL DEFAULT GETDATE(),
    updated_at  DATETIME2     NOT NULL DEFAULT GETDATE()
);

-- Índices para búsquedas frecuentes
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_is_active ON users(is_active);

PRINT 'Tabla users creada exitosamente';
