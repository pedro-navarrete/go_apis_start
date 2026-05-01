package middleware

import (
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/pedro-navarrete/go_apis_start/internal/utils/response"
)

// BasicAuth middleware para autenticación básica (HTTP Basic Auth)
func BasicAuth(users map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Header("WWW-Authenticate", `Basic realm="API"`)
			response.Unauthorized(c, "credenciales requeridas")
			c.Abort()
			return
		}

		// Verificar formato "Basic <base64>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			response.Unauthorized(c, "formato de autenticación inválido")
			c.Abort()
			return
		}

		// Decodificar credenciales Base64
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			response.Unauthorized(c, "credenciales mal formadas")
			c.Abort()
			return
		}

		// Separar username:password
		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 {
			response.Unauthorized(c, "formato de credenciales inválido")
			c.Abort()
			return
		}

		username, pass := credentials[0], credentials[1]

		// Verificar credenciales contra la configuración
		expectedPass, ok := users[username]
		if !ok || expectedPass != pass {
			response.Unauthorized(c, "usuario o contraseña incorrectos")
			c.Abort()
			return
		}

		// Agregar username al contexto
		c.Set("basic_auth_user", username)
		c.Next()
	}
}
