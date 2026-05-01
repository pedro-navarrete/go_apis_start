// Package config maneja la configuración centralizada de la aplicación
package config

import (
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	Server    ServerConfig
	SQLServer SQLServerConfig
	MongoDB   MongoDBConfig
	JWT       JWTConfig
	BasicAuth BasicAuthConfig
	Log       LogConfig
	DBTypes   DBTypesConfig
}

// ServerConfig configuración del servidor HTTP
type ServerConfig struct {
	Port        string
	Environment string
}

// SQLServerConfig configuración de SQL Server
type SQLServerConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// MongoDBConfig configuración de MongoDB
type MongoDBConfig struct {
	URI      string
	Database string
}

// JWTConfig configuración de JWT
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

// BasicAuthConfig configuración de Basic Auth
type BasicAuthConfig struct {
	Users map[string]string // username -> password
}

// LogConfig configuración del logger
type LogConfig struct {
	Level string
}

// DBTypesConfig define qué base de datos usa cada módulo
type DBTypesConfig struct {
	UserDBType    string // sqlserver o mongodb
	ProductDBType string // sqlserver o mongodb
}

// Load carga la configuración desde el archivo .env y variables de entorno
func Load() (*Config, error) {
	// Cargar archivo .env si existe
	_ = godotenv.Load()

	// Configurar Viper
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Valores por defecto
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("SQLSERVER_HOST", "localhost")
	viper.SetDefault("SQLSERVER_PORT", "1433")
	viper.SetDefault("SQLSERVER_USER", "sa")
	viper.SetDefault("SQLSERVER_DB", "go_apis_db")
	viper.SetDefault("MONGODB_URI", "mongodb://localhost:27017")
	viper.SetDefault("MONGODB_DATABASE", "go_apis_db")
	viper.SetDefault("JWT_EXPIRATION", "24h")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("USER_DB_TYPE", "sqlserver")
	viper.SetDefault("PRODUCT_DB_TYPE", "sqlserver")

	// Parsear duración JWT
	jwtExpiration, err := time.ParseDuration(viper.GetString("JWT_EXPIRATION"))
	if err != nil {
		jwtExpiration = 24 * time.Hour
	}

	// Parsear usuarios de Basic Auth
	basicAuthUsers := parseBasicAuthUsers(viper.GetString("BASIC_AUTH_USERS"))

	cfg := &Config{
		Server: ServerConfig{
			Port:        viper.GetString("SERVER_PORT"),
			Environment: viper.GetString("ENVIRONMENT"),
		},
		SQLServer: SQLServerConfig{
			Host:     viper.GetString("SQLSERVER_HOST"),
			Port:     viper.GetString("SQLSERVER_PORT"),
			User:     viper.GetString("SQLSERVER_USER"),
			Password: viper.GetString("SQLSERVER_PASSWORD"),
			Database: viper.GetString("SQLSERVER_DB"),
		},
		MongoDB: MongoDBConfig{
			URI:      viper.GetString("MONGODB_URI"),
			Database: viper.GetString("MONGODB_DATABASE"),
		},
		JWT: JWTConfig{
			Secret:     viper.GetString("JWT_SECRET"),
			Expiration: jwtExpiration,
		},
		BasicAuth: BasicAuthConfig{
			Users: basicAuthUsers,
		},
		Log: LogConfig{
			Level: viper.GetString("LOG_LEVEL"),
		},
		DBTypes: DBTypesConfig{
			UserDBType:    viper.GetString("USER_DB_TYPE"),
			ProductDBType: viper.GetString("PRODUCT_DB_TYPE"),
		},
	}

	return cfg, nil
}

// parseBasicAuthUsers parsea el string "user1:pass1,user2:pass2" en un mapa
func parseBasicAuthUsers(usersStr string) map[string]string {
	users := make(map[string]string)
	if usersStr == "" {
		return users
	}

	pairs := strings.Split(usersStr, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
		if len(parts) == 2 {
			users[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return users
}
