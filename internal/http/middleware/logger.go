package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	applogger "github.com/pedro-navarrete/go_apis_start/internal/utils/logger"
)

// Logger middleware para registrar las solicitudes HTTP
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Tiempo de inicio
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Procesar la solicitud
		c.Next()

		// Calcular duración
		duration := time.Since(start)

		// Construir el path completo
		if query != "" {
			path = path + "?" + query
		}

		// Registrar la solicitud
		applogger.Info("request HTTP",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)
	}
}
