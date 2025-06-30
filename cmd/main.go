package main

import (
	"farma-santi_backend/internal/adapter/database"
	"farma-santi_backend/internal/adapter/handler"
	"farma-santi_backend/internal/adapter/repository"
	"farma-santi_backend/internal/core/service"
	"farma-santi_backend/internal/server"
)

var pool = database.GetDB()

func main() {

	// Inicializar repositorios
	usuarioRepository := repository.NewUsuarioRepository(pool)
	rolRepository := repository.NewRolRepository(pool)
	categoriaRepository := repository.NewCategoriaRepository(pool)
	proveedorRepository := repository.NewProveedorRepository(pool)
	laboratorioRepository := repository.NewLaboratorioRepository(pool)
	productoRepository := repository.NewProductoRepository(pool)
	loteProductoRepository := repository.NewLoteProductoRepository(pool)
	principioActivoRepository := repository.NewPrincipioActivoRepository(pool)
	compraRepository := repository.NewCompraRepository(pool)
	clienteRepository := repository.NewClienteRepository(pool)

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
	compraService := service.NewCompraService(compraRepository)
	clienteService := service.NewClienteService(clienteRepository)

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
	compraHandler := handler.NewCompraHandler(compraService)
	clienteHandler := handler.NewClienteHandler(clienteService)
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
		compraHandler,
		clienteHandler,
	)

	// Iniciar el servidor
	httpServer.Initialize()
}
