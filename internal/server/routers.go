package server

import (
	"farma-santi_backend/internal/core/util"
	"farma-santi_backend/internal/server/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

var limite = limiter.New(limiter.Config{
	Max:        100, // máximo de peticiones permitidas por IP
	Expiration: 5 * time.Minute,
	// Lógica para ralentizar las peticiones
	LimitReached: func(c *fiber.Ctx) error {
		// Si el límite se alcanza, se agrega un retraso
		delay := time.Second * 1
		time.Sleep(delay)
		// Después de la espera, responder con un mensaje de throttling
		return c.Next()
	},
})

func (s *Server) initEndPoints(app *fiber.App) {
	app.Static("/", "./public", fiber.Static{
		ModifyResponse: util.EnviarArchivo,
	})
	api := app.Group("/api") // Crear un grupo para las rutas de la API
	// Inicializar los endpoints de la API
	s.endPointsAPI(api)
}

func (s *Server) endPointsAPI(api fiber.Router) {

	v1 := api.Group("/v1") // Versión 1 de la API
	v1Auth := v1.Group("/auth")
	v1Roles := v1.Group("/roles")
	v1Usuarios := v1.Group("/usuarios")
	v1Categorias := v1.Group("/categorias")
	v1Proveedores := v1.Group("/proveedores")
	v1Laboratorios := v1.Group("/laboratorios")
	v1Productos := v1.Group("/productos")
	v1UsuariosMe := v1Usuarios.Group("/me")
	v1LotesProductos := v1.Group("/lotes-productos")
	v1PrincipiosActivos := v1.Group("/principios-activos")

	// path: /api/v1/usuarios/me
	v1UsuariosMe.Get("", limite, s.usuarioHandler.ObtenerUsuarioActual)

	// Middleware para endpoint
	v1Productos.Use(middleware.HostnameMiddleware, middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"))
	v1Roles.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE"))
	v1Usuarios.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE"))
	v1Categorias.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"))
	v1Proveedores.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE"))
	v1Productos.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"))
	v1Laboratorios.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"))
	v1LotesProductos.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"))
	v1PrincipiosActivos.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"))
	// path: /api/v1/auth
	v1Auth.Post("/login", limite, s.authHandler.Login)
	v1Auth.Get("/logout", limite, s.authHandler.Logout)
	v1Auth.Get("/refresh", limite, s.authHandler.RefreshOrVerify)
	v1Auth.Get("/verify", limite, s.authHandler.RefreshOrVerify)

	// path: /api/v1/roles
	v1Roles.Get("", limite, s.rolHandler.ListarRoles)
	v1Roles.Get("/:rolId", limite, s.rolHandler.ObtenerRolById)
	v1Roles.Post("", limite, s.rolHandler.RegistrarRol)
	v1Roles.Patch("/estado/habilitar/:rolId", limite, s.rolHandler.HabilitarRol)
	v1Roles.Patch("/estado/deshabilitar/:rolId", limite, s.rolHandler.DeshabilitarRol)
	v1Roles.Put("/:rolId", limite, s.rolHandler.ModificarRol)

	// path: /api/v1/usuarios
	v1Usuarios.Get("", limite, s.usuarioHandler.ListarUsuarios)
	v1Usuarios.Get("/:usuarioId", limite, s.usuarioHandler.ObtenerUsuarioDetalle)
	v1Usuarios.Post("", limite, s.usuarioHandler.RegistrarUsuario)
	v1Usuarios.Patch("/estado/habilitar/:usuarioId", limite, s.usuarioHandler.HabilitarUsuarioById)
	v1Usuarios.Patch("/estado/deshabilitar/:usuarioId", limite, s.usuarioHandler.DeshabilitarUsuarioById)
	v1Usuarios.Patch("/password/restablecer/:usuarioId", limite, s.usuarioHandler.RestablecerPassword)
	v1Usuarios.Put("/:usuarioId", limite, s.usuarioHandler.ModificarUsuario)

	//path: /api/v1/categorias
	v1Categorias.Get("/activos", limite, s.categoriaHandler.ListarCategoriasDisponibles)
	v1Categorias.Get("", limite, s.categoriaHandler.ListarCategorias)
	v1Categorias.Get("/:categoriaId", limite, s.categoriaHandler.ObtenerCategoriaById)
	v1Categorias.Post("", limite, s.categoriaHandler.RegistrarCategoria)
	v1Categorias.Put("/:categoriaId", limite, s.categoriaHandler.ModificarCategoria)
	v1Categorias.Patch("/estado/habilitar/:categoriaId", limite, s.categoriaHandler.HabilitarCategoria)
	v1Categorias.Patch("/estado/deshabilitar/:categoriaId", limite, s.categoriaHandler.DeshabilitarCategoria)

	//path: /api/v1/proveedores
	v1Proveedores.Get("", limite, s.proveedorHandler.ListarProveedores)
	v1Proveedores.Get("/:proveedorId", limite, s.proveedorHandler.ObtenerProveedorById)
	v1Proveedores.Post("", limite, s.proveedorHandler.RegistrarProveedor)
	v1Proveedores.Put("/:proveedorId", limite, s.proveedorHandler.ModificarProveedor)
	v1Proveedores.Patch("/estado/habilitar/:proveedorId", limite, s.proveedorHandler.HabilitarProveedor)
	v1Proveedores.Patch("/estado/deshabilitar/:proveedorId", limite, s.proveedorHandler.DeshabilitarProveedor)

	//path: /api/v1/laboratorios
	v1Laboratorios.Get("/activos", limite, s.laboratorioHandler.ListarLaboratoriosDisponibles)
	v1Laboratorios.Get("", limite, s.laboratorioHandler.ListarLaboratorios)
	v1Laboratorios.Get("/:laboratorioId", limite, s.laboratorioHandler.ObtenerLaboratorioById)
	v1Laboratorios.Post("", limite, s.laboratorioHandler.RegistrarLaboratorio)
	v1Laboratorios.Put("/:laboratorioId", limite, s.laboratorioHandler.ModificarLaboratorio)
	v1Laboratorios.Patch("/estado/habilitar/:laboratorioId", limite, s.laboratorioHandler.HabilitarLaboratorio)
	v1Laboratorios.Patch("/estado/deshabilitar/:laboratorioId", limite, s.laboratorioHandler.DeshabilitarLaboratorio)

	//path: /api/v1/productos
	v1Productos.Get("/unidades-medida", limite, s.productoHandler.ListarUnidadesMedida)
	v1Productos.Get("/formas-farmaceuticas", limite, s.productoHandler.ListarFormasFarmaceuticas)
	v1Productos.Get("", limite, s.productoHandler.ListarProductos)
	v1Productos.Get("/:productoId", limite, s.productoHandler.ObtenerProductoById)
	v1Productos.Post("", limite, s.productoHandler.RegistrarProducto)
	v1Productos.Put("/:productoId", limite, s.productoHandler.ModificarProducto)
	v1Productos.Patch("/estado/habilitar/:productoId", limite, s.productoHandler.HabilitarProducto)
	v1Productos.Patch("/estado/deshabilitar/:productoId", limite, s.productoHandler.DeshabilitarProducto)

	//path: /api/v1/lotes-productos
	v1LotesProductos.Get("", limite, s.loteProductoHandler.ListarLotesProductos)
	v1LotesProductos.Get("/:loteProductoId", limite, s.loteProductoHandler.ObtenerLoteProductoById)
	v1LotesProductos.Post("", limite, s.loteProductoHandler.RegistrarLoteProducto)
	v1LotesProductos.Put("/:loteProductoId", limite, s.loteProductoHandler.ModificarLoteProducto)

	//path: /api/v1/principios-activos
	v1PrincipiosActivos.Post("", limite, s.principioActivoHandler.RegistrarPrincipioActivo)
	v1PrincipiosActivos.Put("/:principioActivoId", limite, s.principioActivoHandler.ModificarPrincipioActivo)
	v1PrincipiosActivos.Get("", limite, s.principioActivoHandler.ListarPrincipioActivo)
	v1PrincipiosActivos.Get("/:principioActivoId", limite, s.principioActivoHandler.ObtenerPrincipioActivoById)
}
