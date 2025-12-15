package server

import (
	"farma-santi_backend/internal/core/util"
	"farma-santi_backend/internal/server/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
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

func limited(max int, expiration, delay time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: expiration,
		LimitReached: func(c *fiber.Ctx) error {
			time.Sleep(delay)
			return c.Next()
		},
	})

}
func (s *Server) initEndPoints(app *fiber.App) {
	app.Static("/", "./public", fiber.Static{
		ModifyResponse: util.EnviarArchivo,
	})
	api := app.Group("/api") // Crear un grupo para las rutas de la API
	// Inicializar los endpoints de la API
	s.endPointsShared(api)
	s.endPointsAPI(api)

}

func (s *Server) endPointsAPI(api fiber.Router) {

	v1 := api.Group("/v1") // Versión 1 de la API
	v1Auth := v1.Group("/auth")
	v1Roles := v1.Group("/roles")
	v1Usuarios := v1.Group("/usuarios")
	v1Categorias := v1.Group("/categorias")
	//v1Proveedores := v1.Group("/proveedores")
	v1Laboratorios := v1.Group("/laboratorios")
	v1Productos := v1.Group("/productos")
	v1UsuariosMe := v1Usuarios.Group("/me")
	v1LotesProductos := v1.Group("/lotes-productos")
	v1PrincipiosActivos := v1.Group("/principios-activos")
	v1Compras := v1.Group("/compras")
	v1Clientes := v1.Group("/clientes")
	v1Ventas := v1.Group("/ventas")
	v1Movimientos := v1.Group("/movimientos")
	v1Reportes := v1.Group("/reportes")
	// path: /api/v1/usuarios/me
	v1UsuariosMe.Get("", limited(20, 5*time.Minute, 5*time.Second), s.handlers.Usuario.ObtenerUsuarioActual)

	// Middleware para endpoint
	v1Roles.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE"))
	v1Usuarios.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE"))
	v1Categorias.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"))

	v1PrincipiosActivos.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"))
	v1Clientes.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "FARMACEUTICO"))

	// path: /api/v1/auth
	v1Auth.Post("/login", limited(5, 5*time.Minute, 5*time.Minute), s.handlers.Auth.Login)
	v1Auth.Get("/logout", limited(30, 5*time.Minute, 5*time.Second), s.handlers.Auth.Logout)
	v1Auth.Get("/refresh", limited(50, 5*time.Minute, 5*time.Second), s.handlers.Auth.RefreshOrVerify)
	v1Auth.Get("/verify", limited(50, 5*time.Minute, 5*time.Second), s.handlers.Auth.RefreshOrVerify)

	// path: /api/v1/roles
	v1Roles.Get("", limite, s.handlers.Rol.ListarRoles)
	v1Roles.Get("/:rolId", limite, s.handlers.Rol.ObtenerRolById)
	v1Roles.Post("", limite, s.handlers.Rol.RegistrarRol)
	v1Roles.Patch("/estado/habilitar/:rolId", limite, s.handlers.Rol.HabilitarRol)
	v1Roles.Patch("/estado/deshabilitar/:rolId", limite, s.handlers.Rol.DeshabilitarRol)
	v1Roles.Put("/:rolId", limite, s.handlers.Rol.ModificarRol)

	// path: /api/v1/usuarios
	v1Usuarios.Get("", limite, s.handlers.Usuario.ListarUsuarios)
	v1Usuarios.Get("/:usuarioId", limite, s.handlers.Usuario.ObtenerUsuarioDetalle)
	v1Usuarios.Post("", limite, s.handlers.Usuario.RegistrarUsuario)
	v1Usuarios.Patch("/estado/habilitar/:usuarioId", limite, s.handlers.Usuario.HabilitarUsuarioById)
	v1Usuarios.Patch("/estado/deshabilitar/:usuarioId", limite, s.handlers.Usuario.DeshabilitarUsuarioById)
	v1Usuarios.Patch("/password/restablecer/:usuarioId", limite, s.handlers.Usuario.RestablecerPassword)
	v1Usuarios.Put("/:usuarioId", limite, s.handlers.Usuario.ModificarUsuario)

	//path: /api/v1/categorias
	v1Categorias.Get("/activos", limite, s.handlers.Categoria.ListarCategoriasDisponibles)
	v1Categorias.Get("", limite, s.handlers.Categoria.ListarCategorias)
	v1Categorias.Get("/:categoriaId", limite, s.handlers.Categoria.ObtenerCategoriaById)
	v1Categorias.Post("", limite, s.handlers.Categoria.RegistrarCategoria)
	v1Categorias.Put("/:categoriaId", limite, s.handlers.Categoria.ModificarCategoria)
	v1Categorias.Patch("/estado/habilitar/:categoriaId", limite, s.handlers.Categoria.HabilitarCategoria)
	v1Categorias.Patch("/estado/deshabilitar/:categoriaId", limite, s.handlers.Categoria.DeshabilitarCategoria)

	//path: /api/v1/proveedores
	//v1Proveedores.Get("", limite, s.handlers.Proveedor.ListarProveedores)
	//v1Proveedores.Get("/:proveedorId", limite, s.handlers.Proveedor.ObtenerProveedorById)
	//
	//v1Proveedores.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE"))
	//
	//v1Proveedores.Post("", limite, s.handlers.Proveedor.RegistrarProveedor)
	//v1Proveedores.Put("/:proveedorId", limite, s.handlers.Proveedor.ModificarProveedor)
	//v1Proveedores.Patch("/estado/habilitar/:proveedorId", limite, s.handlers.Proveedor.HabilitarProveedor)
	//v1Proveedores.Patch("/estado/deshabilitar/:proveedorId", limite, s.handlers.Proveedor.DeshabilitarProveedor)

	//path: /api/v1/laboratorios
	v1Laboratorios.Use(middleware.VerifyUserAdminMiddleware)
	v1Laboratorios.Get("/activos", limite, s.handlers.Laboratorio.ListarLaboratoriosDisponibles)
	v1Laboratorios.Get("", limite, s.handlers.Laboratorio.ListarLaboratorios)
	v1Laboratorios.Get("/:laboratorioId", limite, s.handlers.Laboratorio.ObtenerLaboratorioById)
	v1Laboratorios.Use(middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"))
	v1Laboratorios.Post("", limite, s.handlers.Laboratorio.RegistrarLaboratorio)
	v1Laboratorios.Put("/:laboratorioId", limite, s.handlers.Laboratorio.ModificarLaboratorio)
	v1Laboratorios.Patch("/estado/habilitar/:laboratorioId", limite, s.handlers.Laboratorio.HabilitarLaboratorio)
	v1Laboratorios.Patch("/estado/deshabilitar/:laboratorioId", limite, s.handlers.Laboratorio.DeshabilitarLaboratorio)

	//path: /api/v1/productos
	v1Productos.Use(middleware.HostnameMiddleware, middleware.VerifyUserAdminMiddleware)
	v1Productos.Get("/unidades-medida", limite, s.handlers.Producto.ListarUnidadesMedida)
	v1Productos.Get("/formas-farmaceuticas", limite, s.handlers.Producto.ListarFormasFarmaceuticas)
	v1Productos.Get("", limite, s.handlers.Producto.ObtenerListaProductos)
	v1Productos.Get("/:productoId", limite, s.handlers.Producto.ObtenerProductoById)
	v1Productos.Post("", middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"), limite, s.handlers.Producto.RegistrarProducto)
	v1Productos.Put("/:productoId", middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"), limite, s.handlers.Producto.ModificarProducto)
	v1Productos.Patch("/estado/habilitar/:productoId", middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"), limite, s.handlers.Producto.HabilitarProducto)
	v1Productos.Patch("/estado/deshabilitar/:productoId", middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"), limite, s.handlers.Producto.DeshabilitarProducto)

	//path: /api/v1/lotes-productos
	v1LotesProductos.Use(middleware.VerifyUserAdminMiddleware)
	v1LotesProductos.Get("", middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN", "FARMACEUTICO"), limite, s.handlers.LoteProducto.ObtenerListaLotesProductos)
	v1LotesProductos.Get("/byProducto/:productoId", middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"), limite, s.handlers.LoteProducto.ListarLotesProductosByProductoId)
	v1LotesProductos.Get("/:loteProductoId", middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"), limite, s.handlers.LoteProducto.ObtenerLoteProductoById)
	v1LotesProductos.Post("", middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"), limite, s.handlers.LoteProducto.RegistrarLoteProducto)
	v1LotesProductos.Put("/:loteProductoId", middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"), limite, s.handlers.LoteProducto.ModificarLoteProducto)

	//path: /api/v1/principios-activos
	v1PrincipiosActivos.Post("", limite, s.handlers.PrincipioActivo.RegistrarPrincipioActivo)
	v1PrincipiosActivos.Put("/:principioActivoId", limite, s.handlers.PrincipioActivo.ModificarPrincipioActivo)
	v1PrincipiosActivos.Get("", limite, s.handlers.PrincipioActivo.ListarPrincipioActivo)
	v1PrincipiosActivos.Get("/:principioActivoId", limite, s.handlers.PrincipioActivo.ObtenerPrincipioActivoById)

	v1Compras.Use(limite, middleware.VerifyUserAdminMiddleware)
	//path: /api/v1/compras
	v1Compras.Get("", middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN", "FARMACEUTICO"), s.handlers.Compra.ObtenerListaCompras)
	v1Compras.Get("/:compraId", middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN", "FARMACEUTICO"), s.handlers.Compra.ObtenerCompraById)
	v1Compras.Post("", middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"), s.handlers.Compra.RegistrarOrdenCompra)
	v1Compras.Patch("/completar/:compraId", middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"), s.handlers.Compra.RegistrarCompra)
	v1Compras.Put("/:compraId", middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"), s.handlers.Compra.ModificarOrdenCompra)
	v1Compras.Patch("/anular/:compraId", middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"), s.handlers.Compra.AnularOrdenCompra)

	//path: /api/v1/clientes
	v1Clientes.Get("", limite, s.handlers.Cliente.ObtenerListaClientes)
	v1Clientes.Get("/:clienteId", limite, s.handlers.Cliente.ObtenerClienteById)
	v1Clientes.Post("", limite, s.handlers.Cliente.RegistrarCliente)
	v1Clientes.Put("/:clienteId", limite, s.handlers.Cliente.ModificarClienteById)
	v1Clientes.Patch("/estado/habilitar/:clienteId", limite, s.handlers.Cliente.HabilitarCliente)
	v1Clientes.Patch("/estado/deshabilitar/:clienteId", limite, s.handlers.Cliente.DeshabilitarCliente)

	v1Presentaciones := v1.Group("/presentaciones")
	v1Presentaciones.Use(middleware.VerifyUserAdminMiddleware)
	v1Presentaciones.Get("", s.handlers.Presentacion.ObtenerListaPresentaciones)

	//path: /api/v1/ventas
	v1Ventas.Use(middleware.VerifyUserAdminMiddleware, limite)
	v1Ventas.Get("", middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN", "FARMACEUTICO"), s.handlers.Venta.ObtenerListaVentas)
	v1Ventas.Get("/:ventaId", middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "FARMACEUTICO"), s.handlers.Venta.ObtenerVentaById)
	v1Ventas.Post("/registrar", middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "FARMACEUTICO"), s.handlers.Venta.RegistrarVenta)
	v1Ventas.Patch("/anular/:ventaId", middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "FARMACEUTICO"), s.handlers.Venta.AnularVentaById)
	//v1Ventas.Post("/facturar/:ventaId", middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "FARMACEUTICO"), s.handlers.Venta.FacturarVentaById)

	//path: /api/v1/movimientos
	v1Movimientos.Get("", limite, s.handlers.Movimiento.ObtenerListaMovimientos)
	v1Movimientos.Get("/kardex", limite, s.handlers.Movimiento.ObtenerMovimientosKardex)

	//path: /api/stats
	v1Stats := v1.Group("/stats")
	v1Stats.Use(middleware.HostnameMiddleware, middleware.VerifyUserAdminMiddleware, limite)
	v1Stats.Get("/top10Productos", s.handlers.Stat.ObtenerTopProductosVendidos)
	v1Stats.Get("/dashboard", s.handlers.Stat.ObtenerEstadisticasDashboard)

	v1Backups := v1.Group("/backups")
	v1Backups.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE"), limite)
	v1Backups.Get("", s.handlers.Backup.ListarBackups)
	v1Backups.Get("/generate", s.handlers.Backup.DownloadBackup)
	v1Backups.Get("/download/:filename", s.handlers.Backup.DownloadBackupFile)

	//path: /api/v1/reportes
	v1Reportes.Use(middleware.HostnameMiddleware, middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN", "FARMACEUTICO"))
	v1Reportes.Get("/usuarios", s.handlers.Reporte.ReporteUsuariosPDF)
	v1Reportes.Get("/clientes", s.handlers.Reporte.ReporteClientesPDF)
	v1Reportes.Get("/lotes-productos", s.handlers.Reporte.ReporteLotesProductosPDF)
	v1Reportes.Get("/compras", s.handlers.Reporte.ReporteComprasPDF)
	v1Reportes.Get("/compras/:compraId", s.handlers.Reporte.ReporteComprasDetallePDF)
	v1Reportes.Get("/ventas", s.handlers.Reporte.ReporteVentasPDF)
	v1Reportes.Get("/inventario", s.handlers.Reporte.ReporteInventarioPDF)
	v1Reportes.Get("/movimientos", s.handlers.Reporte.ReporteMovimientosPDF)
	v1Reportes.Get("/kardex/:productoId", s.handlers.Reporte.ReporteKardexProductoPDF)
}

func (s *Server) endPointsShared(api fiber.Router) {
	apiShared := api.Group("/shared") // Versión 1 de la API
	v1Productos := apiShared.Group("/productos")
	v1Productos.Use(middleware.HostnameMiddleware)
	v1Productos.Get("", limite, s.handlers.Producto.ObtenerListaProductosShared)
	v1Productos.Get("/formas-farmaceuticas", limite, s.handlers.Producto.ListarFormasFarmaceuticas)
	v1Productos.Get("/:productoId", limite, s.handlers.Producto.ObtenerProductoByIdShared)

	v1Categorias := apiShared.Group("/categorias")
	v1Categorias.Get("", limite, s.handlers.Categoria.ListarCategoriasDisponibles)

	v1Laboratorios := apiShared.Group("/laboratorios")
	v1Laboratorios.Get("", limite, s.handlers.Laboratorio.ListarLaboratoriosDisponibles)

	v1Auth := apiShared.Group("/auth")
	v1Auth.Post("/google/login", limite, s.handlers.Auth.LoginWithGoogle)
	v1Auth.Post("/email/login", limite, s.handlers.Auth.LoginWithEmail)
	v1Auth.Post("/email/register", limite, s.handlers.Auth.RegisterWithEmail)

	v1MisCompras := apiShared.Group("/mis-compras")
	v1MisCompras.Use(middleware.VerifyUsuarioShared)
	v1MisCompras.Get("", limite, s.handlers.Venta.ObtenerListaVentasShared)
	v1MisCompras.Get("/:ventaId", limite, s.handlers.Venta.ObtenerVentaByIdShared)
}
