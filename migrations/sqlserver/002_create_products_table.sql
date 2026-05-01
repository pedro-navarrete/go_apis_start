-- Migración 002: Crear tabla de productos
-- Fecha: 2024-01-01
-- Descripción: Crea la tabla de productos del catálogo

CREATE TABLE products (
    id          NVARCHAR(36)   NOT NULL PRIMARY KEY,
    name        NVARCHAR(255)  NOT NULL,
    description NVARCHAR(MAX)  NULL,
    price       DECIMAL(10,2)  NOT NULL,
    stock       INT            NOT NULL DEFAULT 0,
    is_active   BIT            NOT NULL DEFAULT 1,
    created_at  DATETIME2      NOT NULL DEFAULT GETDATE(),
    updated_at  DATETIME2      NOT NULL DEFAULT GETDATE()
);

-- Índices para búsquedas frecuentes
CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_is_active ON products(is_active);
CREATE INDEX idx_products_price ON products(price);

PRINT 'Tabla products creada exitosamente';
