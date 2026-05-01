# Documentación de la API

## URL Base

```
http://localhost:8080/api
```

---

## Autenticación

### JWT Bearer Token

La mayoría de los endpoints protegidos requieren un token JWT en el header:

```
Authorization: Bearer <token>
```

**Obtener token:**

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "secret123"
  }'
```

**Respuesta:**

```json
{
  "success": true,
  "message": "login exitoso",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "uuid",
      "username": "johndoe",
      "email": "john@example.com",
      "full_name": "John Doe",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

### HTTP Basic Auth

Algunos endpoints usan Basic Auth. Enviar credenciales en Base64:

```
Authorization: Basic <base64(usuario:contraseña)>
```

```bash
curl -u "admin:admin123" http://localhost:8080/api/products \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","price":999.99,"stock":5}'
```

---

## Health Check

### GET /health

Verifica que el servicio está funcionando.

**Respuesta 200:**

```json
{
  "status": "ok",
  "timestamp": "2024-01-01T00:00:00Z",
  "version": "1.0.0"
}
```

```bash
curl http://localhost:8080/api/health
```

---

## Usuarios

### POST /users — Crear usuario

**Body:**

```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "secret123",
  "full_name": "John Doe"
}
```

**Respuesta 201:**

```json
{
  "success": true,
  "message": "usuario creado exitosamente",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"username":"johndoe","email":"john@example.com","password":"secret123","full_name":"John Doe"}'
```

---

### POST /auth/login — Login

**Body:**

```json
{
  "username": "johndoe",
  "password": "secret123"
}
```

**Respuesta 200:**

```json
{
  "success": true,
  "message": "login exitoso",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": { ... }
  }
}
```

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"johndoe","password":"secret123"}'
```

---

### GET /users — Listar usuarios 🔐 JWT

**Query params:** `limit` (default: 10), `offset` (default: 0)

**Respuesta 200:**

```json
{
  "success": true,
  "data": [ ... ],
  "meta": {
    "total": 100,
    "limit": 10,
    "offset": 0,
    "page": 1
  }
}
```

```bash
curl http://localhost:8080/api/users?limit=10&offset=0 \
  -H "Authorization: Bearer <token>"
```

---

### GET /users/:id — Obtener usuario 🔐 JWT

**Respuesta 200:**

```json
{
  "success": true,
  "message": "usuario obtenido",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
    "email": "john@example.com",
    "full_name": "John Doe",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

```bash
curl http://localhost:8080/api/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer <token>"
```

---

### PUT /users/:id — Actualizar usuario 🔐 JWT

**Body (campos opcionales):**

```json
{
  "email": "newemail@example.com",
  "full_name": "John Updated"
}
```

```bash
curl -X PUT http://localhost:8080/api/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"full_name":"John Updated"}'
```

---

### DELETE /users/:id — Eliminar usuario 🔐 JWT

**Respuesta 204:** Sin contenido.

```bash
curl -X DELETE http://localhost:8080/api/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer <token>"
```

---

## Productos

### GET /products — Listar productos

**Query params:** `limit` (default: 10), `offset` (default: 0)

**Respuesta 200:**

```json
{
  "success": true,
  "data": [ ... ],
  "meta": {
    "total": 50,
    "limit": 10,
    "offset": 0,
    "page": 1
  }
}
```

```bash
curl http://localhost:8080/api/products?limit=5&offset=0
```

---

### GET /products/:id — Obtener producto

**Respuesta 200:**

```json
{
  "success": true,
  "message": "producto obtenido",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Laptop Pro",
    "description": "Laptop de alta gama",
    "price": 1299.99,
    "stock": 25,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

```bash
curl http://localhost:8080/api/products/550e8400-e29b-41d4-a716-446655440001
```

---

### POST /products — Crear producto 🔑 Basic Auth

**Body:**

```json
{
  "name": "Laptop Pro",
  "description": "Laptop de alta gama",
  "price": 1299.99,
  "stock": 25
}
```

**Respuesta 201:**

```json
{
  "success": true,
  "message": "producto creado exitosamente",
  "data": { ... }
}
```

```bash
curl -X POST http://localhost:8080/api/products \
  -u "admin:admin123" \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop Pro","description":"Laptop de alta gama","price":1299.99,"stock":25}'
```

---

### PUT /products/:id — Actualizar producto 🔐 JWT

**Body (campos opcionales):**

```json
{
  "name": "Laptop Pro 2024",
  "price": 1199.99,
  "stock": 30
}
```

```bash
curl -X PUT http://localhost:8080/api/products/550e8400-e29b-41d4-a716-446655440001 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"price":1199.99,"stock":30}'
```

---

### DELETE /products/:id — Eliminar producto 🔐 JWT

**Respuesta 204:** Sin contenido.

```bash
curl -X DELETE http://localhost:8080/api/products/550e8400-e29b-41d4-a716-446655440001 \
  -H "Authorization: Bearer <token>"
```

---

## Códigos de Error

| Código | Descripción                              |
|--------|------------------------------------------|
| 200    | OK - Operación exitosa                   |
| 201    | Created - Recurso creado                 |
| 204    | No Content - Eliminación exitosa         |
| 400    | Bad Request - Datos inválidos            |
| 401    | Unauthorized - Sin autenticación         |
| 404    | Not Found - Recurso no encontrado        |
| 409    | Conflict - Recurso ya existe             |
| 500    | Internal Server Error - Error del server |

**Formato de error:**

```json
{
  "success": false,
  "message": "descripción del error",
  "error": "detalle técnico del error"
}
```
