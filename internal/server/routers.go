package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

var limite = limiter.New(limiter.Config{
	Max:        200, // máximo de peticiones permitidas por IP
	Expiration: 5 * time.Minute,
	// Lógica para ralentizar las peticiones
	LimitReached: func(c *fiber.Ctx) error {
		// Si el límite se alcanza, se agrega un retraso
		//delay := time.Second * 1
		//time.Sleep(delay)
		// Después de la espera, responder con un mensaje de throttling
		return c.Next()
	},
})

func (s *Server) initEndPoints(app *fiber.App) {
	api := app.Group("/api") // Crear un grupo para las rutas de la API
	// Inicializar los endpoints de la API
	s.endPointsAPI(api)
}

func (s *Server) endPointsAPI(api fiber.Router) {

	v1 := api.Group("/v1") // Versión 1 de la API
	v1Auth := v1.Group("/auth")
	v1Roles := v1.Group("/roles")
	v1Usuario := v1.Group("/usuarios")
	v1Categoria := v1.Group("/categorias")
	v1Proveedor := v1.Group("/proveedores")
	// path: /api/v1/auth
	v1Auth.Post("/login", limite, s.authHandler.Login)
	v1Auth.Get("/logout", limite, s.authHandler.Logout)
	v1Auth.Get("/refresh", limite, s.authHandler.RefreshOrVerify)
	v1Auth.Get("/verify", limite, s.authHandler.RefreshOrVerify)

	// path: /api/v1/roles
	v1Roles.Get("", limite, s.rolHandler.ListarRoles)
	v1Roles.Get("/:rolId", limite, s.rolHandler.ObtenerRolById)
	v1Roles.Post("", limite, s.rolHandler.RegistrarRol)
	v1Roles.Patch("/status/:rolId", limite, s.rolHandler.ModificarEstadoRol)
	v1Roles.Put("/:rolId", limite, s.rolHandler.ModificarRol)

	// path: /api/v1/usuarios
	v1Usuario.Get("", limite, s.usuarioHandler.ListarUsuarios)
	v1Usuario.Get("/:usuarioId", limite, s.usuarioHandler.ObtenerUsuarioDetalle)
	v1Usuario.Post("", limite, s.usuarioHandler.RegistrarUsuario)
	v1Usuario.Patch("/status/:usuarioId", limite, s.usuarioHandler.ModificarEstadoUsuario)
	v1Usuario.Put("/:usuarioId", limite, s.usuarioHandler.ModificarUsuario)

	//path: /api/v1/categorias
	v1Categoria.Get("", limite, s.categoriaHandler.ListarCategorias)
	v1Categoria.Get("/:categoriaId", limite, s.categoriaHandler.ObtenerCategoriaById)
	v1Categoria.Post("", limite, s.categoriaHandler.RegistrarCategoria)
	v1Categoria.Put("/:categoriaId", limite, s.categoriaHandler.ModificarCategoria)
	v1Categoria.Patch("/status/:categoriaId", limite, s.categoriaHandler.ModificarEstadoCategoria)

	//path: /api/v1/proveedores
	v1Proveedor.Get("", limite, s.proveedorHandler.ListarProveedores)
	v1Proveedor.Get("/:proveedorId", limite, s.proveedorHandler.ObtenerProveedorById)
	v1Proveedor.Post("", limite, s.proveedorHandler.RegistrarProveedor)
	v1Proveedor.Put("/:proveedorId", limite, s.proveedorHandler.ModificarProveedor)
	v1Proveedor.Patch("/status/:proveedorId", limite, s.proveedorHandler.ModificarEstadoProveedor)
}
