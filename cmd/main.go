package main

import (
	"farma-santi_backend/internal/adapter/database"
	"farma-santi_backend/internal/adapter/handler"
	"farma-santi_backend/internal/adapter/repository"
	"farma-santi_backend/internal/core/service"
	"farma-santi_backend/internal/logger"
	"farma-santi_backend/internal/server"
)

var db = &database.DBInstance

func main() {
	// Configurar el logger global
	if err := logger.SetupLogger(); err != nil {
		return
	}
	// Usar el logger
	logger.ZapLog.Info("Aplicación iniciada")

	// Asegurarse de cerrar el logger al final
	defer logger.Shutdown()

	// Inicializar repositorios
	usuarioRepository := repository.NewUsuarioRepository(db)
	rolRepository := repository.NewRolRepository(db)
	categoriaRepository := repository.NewCategoriaRepository(db)
	proveedorRepository := repository.NewProveedorRepository(db)
	// Inicializar servicios
	authService := service.NewAuthService(usuarioRepository)
	rolService := service.NewRolService(rolRepository)
	usuarioService := service.NewUsuarioService(usuarioRepository)
	categoriaService := service.NewCategoriaService(categoriaRepository)
	proveedorService := service.NewProveedorService(proveedorRepository)
	// Inicializar manejadores/controladores/handlers
	authHandler := handler.NewAuthHandler(authService)
	rolHandler := handler.NewRolHandler(rolService)
	usuarioHandler := handler.NewUsuarioHandler(usuarioService)
	categoriaHandler := handler.NewCategoriaHandler(categoriaService)
	proveedorHandler := handler.NewProveedorHandler(proveedorService)
	// Inicializar el servidor HTTP
	httpServer := server.NewServer(authHandler, rolHandler, usuarioHandler, categoriaHandler, proveedorHandler)

	// Iniciar el servidor
	httpServer.Initialize()
}
