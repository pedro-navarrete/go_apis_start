package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	applogger "github.com/pedro-navarrete/go_apis_start/internal/utils/logger"
	"github.com/pedro-navarrete/go_apis_start/internal/utils/response"
)

// Recovery middleware para capturar panics y retornar 500
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Registrar el panic
				applogger.Error("panic capturado",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				// Retornar error 500
				response.Error(c, http.StatusInternalServerError, "error interno del servidor", nil)
				c.Abort()
			}
		}()

		c.Next()
	}
}
