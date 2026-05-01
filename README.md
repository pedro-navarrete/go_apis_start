# go_apis_start

![Go Version](https://img.shields.io/badge/Go-1.22-blue)
![Gin](https://img.shields.io/badge/Gin-v1.10-green)
![License](https://img.shields.io/badge/License-MIT-yellow)

API REST en Go con soporte dual para **SQL Server** y **MongoDB**, autenticación JWT y Basic Auth, diseñada con arquitectura orientada al dominio.

## Características

- 🚀 **Framework**: Gin v1.10 (alto rendimiento)
- 🗄️ **SQL Server**: GORM v1.25 con driver SQL Server
- 🍃 **MongoDB**: mongo-driver v1.14 oficial
- 🔐 **JWT**: golang-jwt/jwt v5 para autenticación Bearer
- 🔑 **Basic Auth**: Autenticación básica HTTP configurable
- 📝 **Logger**: Zap estructurado (Uber)
- ✅ **Validación**: go-playground/validator v10
- ⚙️ **Config**: Viper + godotenv
- 🐳 **Docker**: Docker Compose con SQL Server 2022 + MongoDB 7
- 🔄 **Graceful Shutdown**: Cierre limpio del servidor
- 🏛️ **Arquitectura**: Domain-Driven Design con Repository Pattern

## Arquitectura del Proyecto

```
go_apis_start/
├── cmd/
│   └── api/
│       └── main.go                    # Punto de entrada
├── internal/
│   ├── config/
│   │   └── config.go                  # Configuración centralizada
│   ├── domain/
│   │   ├── user/
│   │   │   ├── model.go               # Entidad User + DTOs
│   │   │   ├── repository.go          # Interfaz repositorio
│   │   │   └── service.go             # Lógica de negocio
│   │   └── product/
│   │       ├── model.go               # Entidad Product + DTOs
│   │       ├── repository.go          # Interfaz repositorio
│   │       └── service.go             # Lógica de negocio
│   ├── http/
│   │   ├── handlers/
│   │   │   ├── health_handler.go
│   │   │   ├── user_handler.go
│   │   │   └── product_handler.go
│   │   ├── middleware/
│   │   │   ├── jwt.go                 # Middleware JWT
│   │   │   ├── auth.go                # Middleware Basic Auth
│   │   │   ├── cors.go
│   │   │   ├── logger.go
│   │   │   └── recovery.go
│   │   └── routes/
│   │       └── routes.go              # Registro de rutas
│   ├── infrastructure/
│   │   ├── database/
│   │   │   ├── sqlserver.go           # Conexión SQL Server
│   │   │   └── mongodb.go             # Conexión MongoDB
│   │   └── repository/
│   │       ├── user_repository_sqlserver.go
│   │       ├── user_repository_mongo.go
│   │       ├── product_repository_sqlserver.go
│   │       └── product_repository_mongo.go
│   └── utils/
│       ├── logger/                    # Zap logger
│       ├── response/                  # Respuestas HTTP estandarizadas
│       ├── validator/                 # Validación de structs
│       └── password/                  # bcrypt hash/verify
├── pkg/
│   └── errors/
│       └── errors.go                  # Errores de dominio
├── migrations/
│   ├── sqlserver/                     # Scripts SQL
│   └── mongodb/                       # Scripts de inicialización
├── docs/
│   └── API.md                         # Documentación de la API
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
└── .env.example
```

## Requisitos Previos

- **Go** 1.22 o superior
- **Docker** y **Docker Compose**
- (Opcional) **golangci-lint** para el linter

## Instalación y Configuración

### 1. Clonar el repositorio

```bash
git clone https://github.com/pedro-navarrete/go_apis_start.git
cd go_apis_start
```

### 2. Instalar dependencias

```bash
make install
# o manualmente:
go mod download && go mod tidy
```

### 3. Configurar variables de entorno

```bash
cp .env.example .env
# Editar .env con tus valores
```

## Configuración

Copia `.env.example` a `.env` y ajusta los valores:

```env
SERVER_PORT=8080
ENVIRONMENT=development

SQLSERVER_HOST=localhost
SQLSERVER_PORT=1433
SQLSERVER_USER=sa
SQLSERVER_PASSWORD=YourStrong@Passw0rd
SQLSERVER_DB=go_apis_db

MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=go_apis_db

JWT_SECRET=tu-clave-secreta-super-segura
JWT_EXPIRATION=24h

BASIC_AUTH_USERS=admin:admin123,developer:dev123

LOG_LEVEL=debug

USER_DB_TYPE=sqlserver      # sqlserver o mongodb
PRODUCT_DB_TYPE=sqlserver   # sqlserver o mongodb
```

## Ejecución

### Desarrollo completo (Docker + App)

```bash
make dev
```

### Solo levantar bases de datos

```bash
make docker-up
```

### Ejecutar la aplicación

```bash
make run
```

### Compilar binario

```bash
make build
./bin/api
```

## Endpoints de la API

| Método | Endpoint            | Auth       | Descripción                  |
|--------|---------------------|------------|------------------------------|
| GET    | /api/health         | Ninguna    | Health check                 |
| POST   | /api/auth/login     | Ninguna    | Login → devuelve JWT         |
| POST   | /api/users          | Ninguna    | Crear usuario                |
| GET    | /api/users          | JWT Bearer | Listar usuarios              |
| GET    | /api/users/:id      | JWT Bearer | Obtener usuario por ID       |
| PUT    | /api/users/:id      | JWT Bearer | Actualizar usuario           |
| DELETE | /api/users/:id      | JWT Bearer | Eliminar usuario             |
| GET    | /api/products       | Ninguna    | Listar productos             |
| GET    | /api/products/:id   | Ninguna    | Obtener producto por ID      |
| POST   | /api/products       | Basic Auth | Crear producto               |
| PUT    | /api/products/:id   | JWT Bearer | Actualizar producto          |
| DELETE | /api/products/:id   | JWT Bearer | Eliminar producto            |

## Autenticación

### JWT Bearer Token

```bash
# 1. Obtener token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 2. Usar token en requests
curl http://localhost:8080/api/users \
  -H "Authorization: Bearer <token>"
```

### HTTP Basic Auth

```bash
curl -X POST http://localhost:8080/api/products \
  -u "admin:admin123" \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","price":999.99,"stock":10}'
```

## Bases de Datos

Puedes configurar cada módulo para usar SQL Server o MongoDB de forma independiente:

```env
USER_DB_TYPE=sqlserver    # usuarios en SQL Server
PRODUCT_DB_TYPE=mongodb   # productos en MongoDB
```

## Docker

```bash
make docker-up      # Levantar SQL Server + MongoDB + Mongo Express
make docker-down    # Detener servicios
make docker-logs    # Ver logs
make docker-clean   # Eliminar volúmenes
```

**Servicios disponibles:**
- SQL Server: `localhost:1433`
- MongoDB: `localhost:27017`
- Mongo Express: http://localhost:8081 (admin / admin)

## Makefile

```bash
make help           # Ver todos los comandos
make run            # Ejecutar la app
make build          # Compilar binario
make test           # Ejecutar tests
make clean          # Limpiar archivos generados
make install        # Instalar dependencias
make fmt            # Formatear código
make lint           # Ejecutar golangci-lint
make dev            # Docker up + ejecutar app
make docker-up      # Levantar Docker
make docker-down    # Detener Docker
make docker-logs    # Ver logs Docker
make docker-clean   # Limpiar volúmenes Docker
```
