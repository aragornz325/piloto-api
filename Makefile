# Ruta al archivo main.go que contiene los comentarios principales de Swagger
SWAG_MAIN=cmd/server/main.go

# Comando para generar documentación Swagger
swagger:
	@echo "🔧 Generando documentación Swagger desde $(SWAG_MAIN)..."
	@swag init -g $(SWAG_MAIN) -o ./docs
	@echo "✅ Documentación Swagger generada en ./docs"

# Comando para limpiar la documentación Swagger
swagger-clean:
	@echo "🧹 Eliminando documentación Swagger..."
	@rm -rf ./docs/swagger.json ./docs/swagger.yaml ./docs/docs.go
	@echo "✅ Documentación eliminada"

# Comando completo para limpiar y regenerar
swagger-rebuild: swagger-clean swagger

# Comando para arrancar el servidor
run:
	@echo "🚀 Iniciando el servidor..."
	@go run cmd/server/main.go
	@echo "✅ Servidor iniciado"
	@echo "📜 Documentación disponible en http://localhost:3800/swagger/index.html"