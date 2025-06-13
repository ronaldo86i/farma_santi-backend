package main

import (
	"farma-santi_backend/internal/adapter/database"
	"farma-santi_backend/internal/adapter/handler"
	"farma-santi_backend/internal/adapter/repository"
	"farma-santi_backend/internal/core/service"
	"farma-santi_backend/internal/server"
)

var db = &database.DBInstance

func main() {

	// Inicializar repositorios
	usuarioRepository := repository.NewUsuarioRepository(db)
	rolRepository := repository.NewRolRepository(db)
	categoriaRepository := repository.NewCategoriaRepository(db)
	proveedorRepository := repository.NewProveedorRepository(db)
	laboratorioRepository := repository.NewLaboratorioRepository(db)
	productoRepository := repository.NewProductoRepository(db)
	loteProductoRepository := repository.NewLoteProductoRepository(db)
	principioActivoRepository := repository.NewPrincipioActivoRepository(db)
	// Inicializar servicios
	authService := service.NewAuthService(usuarioRepository)
	rolService := service.NewRolService(rolRepository)
	usuarioService := service.NewUsuarioService(usuarioRepository)
	categoriaService := service.NewCategoriaService(categoriaRepository)
	proveedorService := service.NewProveedorService(proveedorRepository)
	laboratorioService := service.NewLaboratorioService(laboratorioRepository)
	productoService := service.NewProductoService(productoRepository)
	loteProductoService := service.NewLoteProductoService(loteProductoRepository)
	principioActivoService := service.NewPrincipioActivoService(principioActivoRepository)
	// Inicializar manejadores/controladores/handlers
	authHandler := handler.NewAuthHandler(authService)
	rolHandler := handler.NewRolHandler(rolService)
	usuarioHandler := handler.NewUsuarioHandler(usuarioService)
	categoriaHandler := handler.NewCategoriaHandler(categoriaService)
	proveedorHandler := handler.NewProveedorHandler(proveedorService)
	laboratorioHandler := handler.NewLaboratorioHandler(laboratorioService)
	productoHandler := handler.NewProductoHandler(productoService)
	loteProductoHandler := handler.NewLoteProductoHandler(loteProductoService)
	principioActivoHandler := handler.NewPrincipioActivoHandler(principioActivoService)
	// Inicializar el servidor HTTP
	httpServer := server.NewServer(
		authHandler,
		rolHandler,
		usuarioHandler,
		categoriaHandler,
		proveedorHandler,
		laboratorioHandler,
		productoHandler,
		loteProductoHandler,
		principioActivoHandler,
	)

	// Iniciar el servidor
	httpServer.Initialize()
}
