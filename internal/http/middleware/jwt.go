// Package middleware contiene los middlewares de la aplicación
package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/pedro-navarrete/go_apis_start/internal/utils/response"
)

// JWTClaims define los claims del token JWT
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTAuth middleware para validar tokens JWT Bearer
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "token de autorización requerido")
			c.Abort()
			return
		}

		// Verificar formato "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "formato de token inválido, usar: Bearer <token>")
			c.Abort()
			return
		}

		tokenStr := parts[1]

		// Parsear y validar el token
		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			// Verificar el método de firma
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			response.Unauthorized(c, "token inválido o expirado")
			c.Abort()
			return
		}

		// Agregar información del usuario al contexto
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
