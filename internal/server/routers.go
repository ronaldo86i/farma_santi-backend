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
	v1Compras := v1.Group("/compras")
	v1Clientes := v1.Group("/clientes")

	// path: /api/v1/usuarios/me
	v1UsuariosMe.Get("", limited(20, 5*time.Minute, 5*time.Second), s.handlers.Usuario.ObtenerUsuarioActual)

	// Middleware para endpoint
	v1Productos.Use(middleware.HostnameMiddleware, middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"))
	v1Roles.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE"))
	v1Usuarios.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE"))
	v1Categorias.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"))

	v1Productos.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"))
	v1Laboratorios.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"))
	v1LotesProductos.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"))
	v1PrincipiosActivos.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"))
	v1Compras.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE", "AUXILIAR DE ALMACEN"))

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
	v1Proveedores.Get("", limite, s.handlers.Proveedor.ListarProveedores)
	v1Proveedores.Get("/:proveedorId", limite, s.handlers.Proveedor.ObtenerProveedorById)

	v1Proveedores.Use(middleware.VerifyUserAdminMiddleware, middleware.VerifyRolesMiddleware("ADMIN", "GERENTE"))

	v1Proveedores.Post("", limite, s.handlers.Proveedor.RegistrarProveedor)
	v1Proveedores.Put("/:proveedorId", limite, s.handlers.Proveedor.ModificarProveedor)
	v1Proveedores.Patch("/estado/habilitar/:proveedorId", limite, s.handlers.Proveedor.HabilitarProveedor)
	v1Proveedores.Patch("/estado/deshabilitar/:proveedorId", limite, s.handlers.Proveedor.DeshabilitarProveedor)

	//path: /api/v1/laboratorios
	v1Laboratorios.Get("/activos", limite, s.handlers.Laboratorio.ListarLaboratoriosDisponibles)
	v1Laboratorios.Get("", limite, s.handlers.Laboratorio.ListarLaboratorios)
	v1Laboratorios.Get("/:laboratorioId", limite, s.handlers.Laboratorio.ObtenerLaboratorioById)
	v1Laboratorios.Post("", limite, s.handlers.Laboratorio.RegistrarLaboratorio)
	v1Laboratorios.Put("/:laboratorioId", limite, s.handlers.Laboratorio.ModificarLaboratorio)
	v1Laboratorios.Patch("/estado/habilitar/:laboratorioId", limite, s.handlers.Laboratorio.HabilitarLaboratorio)
	v1Laboratorios.Patch("/estado/deshabilitar/:laboratorioId", limite, s.handlers.Laboratorio.DeshabilitarLaboratorio)

	//path: /api/v1/productos
	v1Productos.Get("/unidades-medida", limite, s.handlers.Producto.ListarUnidadesMedida)
	v1Productos.Get("/formas-farmaceuticas", limite, s.handlers.Producto.ListarFormasFarmaceuticas)
	v1Productos.Get("", limite, s.handlers.Producto.ListarProductos)
	v1Productos.Get("/:productoId", limite, s.handlers.Producto.ObtenerProductoById)
	v1Productos.Post("", limite, s.handlers.Producto.RegistrarProducto)
	v1Productos.Put("/:productoId", limite, s.handlers.Producto.ModificarProducto)
	v1Productos.Patch("/estado/habilitar/:productoId", limite, s.handlers.Producto.HabilitarProducto)
	v1Productos.Patch("/estado/deshabilitar/:productoId", limite, s.handlers.Producto.DeshabilitarProducto)

	//path: /api/v1/lotes-productos
	v1LotesProductos.Get("", limite, s.handlers.LoteProducto.ListarLotesProductos)
	v1LotesProductos.Get("/byProducto/:productoId", limite, s.handlers.LoteProducto.ListarLotesProductosByProductoId)
	v1LotesProductos.Get("/:loteProductoId", limite, s.handlers.LoteProducto.ObtenerLoteProductoById)
	v1LotesProductos.Post("", limite, s.handlers.LoteProducto.RegistrarLoteProducto)
	v1LotesProductos.Put("/:loteProductoId", limite, s.handlers.LoteProducto.ModificarLoteProducto)

	//path: /api/v1/principios-activos
	v1PrincipiosActivos.Post("", limite, s.handlers.PrincipioActivo.RegistrarPrincipioActivo)
	v1PrincipiosActivos.Put("/:principioActivoId", limite, s.handlers.PrincipioActivo.ModificarPrincipioActivo)
	v1PrincipiosActivos.Get("", limite, s.handlers.PrincipioActivo.ListarPrincipioActivo)
	v1PrincipiosActivos.Get("/:principioActivoId", limite, s.handlers.PrincipioActivo.ObtenerPrincipioActivoById)

	//path: /api/v1/compras
	v1Compras.Get("", limite, s.handlers.Compra.ObtenerListaCompras)
	v1Compras.Get("/:compraId", limite, s.handlers.Compra.ObtenerCompraById)
	v1Compras.Post("", limite, s.handlers.Compra.RegistrarOrdenCompra)
	v1Compras.Patch("/completar/:compraId", limite, s.handlers.Compra.RegistrarCompra)
	v1Compras.Put("/:compraId", limite, s.handlers.Compra.ModificarOrdenCompra)
	v1Compras.Patch("/anular/:compraId", limite, s.handlers.Compra.AnularOrdenCompra)

	//path: /api/v1/clientes
	v1Clientes.Get("", limite, s.handlers.Cliente.ObtenerListaClientes)
	v1Clientes.Get("/:clienteId", limite, s.handlers.Cliente.ObtenerClienteById)
	v1Clientes.Post("", limite, s.handlers.Cliente.RegistrarCliente)
	v1Clientes.Put("/:clienteId", limite, s.handlers.Cliente.ModificarClienteById)
	v1Clientes.Patch("/estado/habilitar/:clienteId", limite, s.handlers.Cliente.HabilitarCliente)
	v1Clientes.Patch("/estado/deshabilitar/:clienteId", limite, s.handlers.Cliente.DeshabilitarCliente)
}
