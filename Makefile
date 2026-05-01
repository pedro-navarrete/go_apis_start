.PHONY: help run build test clean docker-up docker-down install dev fmt lint

# Variables
BINARY_NAME=api
MAIN_PATH=cmd/api/main.go

help: ## Mostrar esta ayuda
@echo "Comandos disponibles:"
@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: ## Ejecutar la aplicación
@echo "🚀 Ejecutando la aplicación..."
go run $(MAIN_PATH)

build: ## Compilar la aplicación
@echo "🔨 Compilando..."
@mkdir -p bin
go build -o bin/$(BINARY_NAME) $(MAIN_PATH)
@echo "✅ Binario creado en bin/$(BINARY_NAME)"

test: ## Ejecutar tests
@echo "🧪 Ejecutando tests..."
go test -v ./...

clean: ## Limpiar archivos generados
@echo "🧹 Limpiando..."
@rm -rf bin/
@rm -f coverage.out
@echo "✅ Limpieza completada"

install: ## Instalar dependencias
@echo "📦 Instalando dependencias..."
go mod download
go mod tidy
@echo "✅ Dependencias instaladas"

docker-up: ## Levantar servicios de Docker
@echo "🐳 Levantando servicios..."
docker-compose up -d
@echo "✅ Servicios levantados"
@echo "📊 SQL Server: localhost:1433"
@echo "📊 MongoDB: localhost:27017"
@echo "📊 Mongo Express: http://localhost:8081 (admin / admin)"

docker-down: ## Detener servicios de Docker
@echo "🛑 Deteniendo servicios..."
docker-compose down
@echo "✅ Servicios detenidos"

docker-logs: ## Ver logs de Docker
docker-compose logs -f

docker-clean: ## Limpiar volúmenes de Docker
@echo "🧹 Limpiando volúmenes..."
docker-compose down -v
@echo "✅ Volúmenes limpiados"

dev: docker-up ## Levantar Docker y ejecutar la app
@echo "⏳ Esperando que las bases de datos estén listas..."
@sleep 10
@$(MAKE) run

fmt: ## Formatear código
@echo "🎨 Formateando código..."
go fmt ./...
@echo "✅ Código formateado"

lint: ## Ejecutar linter
@echo "🔍 Ejecutando linter..."
golangci-lint run ./...
