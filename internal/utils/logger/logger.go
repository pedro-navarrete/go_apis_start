// Package logger proporciona funciones de logging estructurado usando Zap
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logger instancia global del logger
var logger *zap.Logger

// Init inicializa el logger según el entorno
func Init(level, environment string) error {
	var config zap.Config

	if environment == "production" {
		// Producción: formato JSON
		config = zap.NewProductionConfig()
	} else {
		// Desarrollo: formato legible en consola
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Configurar nivel de log
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	var err error
	logger, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

// Get retorna la instancia del logger
func Get() *zap.Logger {
	if logger == nil {
		// Logger por defecto si no se ha inicializado
		logger, _ = zap.NewDevelopment()
	}
	return logger
}

// Debug registra un mensaje de nivel debug
func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

// Info registra un mensaje de nivel info
func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

// Warn registra un mensaje de nivel warn
func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

// Error registra un mensaje de nivel error
func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

// Fatal registra un mensaje fatal y termina la aplicación
func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

// Sync sincroniza el buffer del logger
func Sync() {
	_ = Get().Sync()
}
