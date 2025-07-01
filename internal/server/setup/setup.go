package setup

import (
	"farma-santi_backend/internal/adapter/handler"
	"farma-santi_backend/internal/adapter/repository"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/service"
	"farma-santi_backend/internal/postgresql"
	"github.com/joho/godotenv"
	"log"
	"sync"
)

type Repository struct {
	Categoria       port.CategoriaRepository
	Cliente         port.ClienteRepository
	Compras         port.CompraRepository
	Laboratorio     port.LaboratorioRepository
	LoteProducto    port.LoteProductoRepository
	PrincipioActivo port.PrincipioActivoRepository
	Producto        port.ProductoRepository
	Proveedor       port.ProveedorRepository
	Rol             port.RolRepository
	Usuario         port.UsuarioRepository
}

type Service struct {
	Auth            port.AuthService
	Categoria       port.CategoriaService
	Cliente         port.ClienteService
	Compras         port.CompraService
	Laboratorio     port.LaboratorioService
	LoteProducto    port.LoteProductoService
	PrincipioActivo port.PrincipioActivoService
	Producto        port.ProductoService
	Proveedor       port.ProveedorService
	Rol             port.RolService
	Usuario         port.UsuarioService
}

type Handler struct {
	Auth            port.AuthHandler
	Categoria       port.CategoriaHandler
	Cliente         port.ClienteHandler
	Compra          port.CompraHandler
	Laboratorio     port.LaboratorioHandler
	LoteProducto    port.LoteProductoHandler
	PrincipioActivo port.PrincipioActivoHandler
	Producto        port.ProductoHandler
	Proveedor       port.ProveedorHandler
	Rol             port.RolHandler
	Usuario         port.UsuarioHandler
}

type Dependencies struct {
	Repository Repository
	Service    Service
	Handler    Handler
}

var (
	instance *Dependencies
	once     sync.Once
)

func GetDependencies() *Dependencies {
	return instance
}

func initEnv(filenames ...string) error {
	err := godotenv.Load(filenames...)
	if err != nil {
		return err
	}
	return nil
}

func initDB() error {
	err := postgresql.Connection()
	if err != nil {
		return err
	}
	return nil
}

func Init() {
	once.Do(func() {
		if err := initEnv(".env"); err != nil {
			log.Fatalf("Fallo al inicializar variables de entorno: %v", err)
		}

		if err := initDB(); err != nil {
			log.Fatalf("Fallo en conectar a la base de datos: %v", err)
		}
		var pool = postgresql.GetDB()
		d := &Dependencies{}

		// Repositories
		d.Repository.Usuario = repository.NewUsuarioRepository(pool)
		d.Repository.Rol = repository.NewRolRepository(pool)
		d.Repository.Categoria = repository.NewCategoriaRepository(pool)
		d.Repository.Proveedor = repository.NewProveedorRepository(pool)
		d.Repository.Laboratorio = repository.NewLaboratorioRepository(pool)
		d.Repository.Producto = repository.NewProductoRepository(pool)
		d.Repository.LoteProducto = repository.NewLoteProductoRepository(pool)
		d.Repository.PrincipioActivo = repository.NewPrincipioActivoRepository(pool)
		d.Repository.Compras = repository.NewCompraRepository(pool)
		d.Repository.Cliente = repository.NewClienteRepository(pool)

		// Services
		d.Service.Auth = service.NewAuthService(d.Repository.Usuario)
		d.Service.Usuario = service.NewUsuarioService(d.Repository.Usuario)
		d.Service.Rol = service.NewRolService(d.Repository.Rol)
		d.Service.Categoria = service.NewCategoriaService(d.Repository.Categoria)
		d.Service.Proveedor = service.NewProveedorService(d.Repository.Proveedor)
		d.Service.Laboratorio = service.NewLaboratorioService(d.Repository.Laboratorio)
		d.Service.Producto = service.NewProductoService(d.Repository.Producto)
		d.Service.LoteProducto = service.NewLoteProductoService(d.Repository.LoteProducto)
		d.Service.PrincipioActivo = service.NewPrincipioActivoService(d.Repository.PrincipioActivo)
		d.Service.Compras = service.NewCompraService(d.Repository.Compras)
		d.Service.Cliente = service.NewClienteService(d.Repository.Cliente)

		// Handlers
		d.Handler.Auth = handler.NewAuthHandler(d.Service.Auth)
		d.Handler.Usuario = handler.NewUsuarioHandler(d.Service.Usuario)
		d.Handler.Rol = handler.NewRolHandler(d.Service.Rol)
		d.Handler.Categoria = handler.NewCategoriaHandler(d.Service.Categoria)
		d.Handler.Proveedor = handler.NewProveedorHandler(d.Service.Proveedor)
		d.Handler.Laboratorio = handler.NewLaboratorioHandler(d.Service.Laboratorio)
		d.Handler.Producto = handler.NewProductoHandler(d.Service.Producto)
		d.Handler.LoteProducto = handler.NewLoteProductoHandler(d.Service.LoteProducto)
		d.Handler.PrincipioActivo = handler.NewPrincipioActivoHandler(d.Service.PrincipioActivo)
		d.Handler.Compra = handler.NewCompraHandler(d.Service.Compras)
		d.Handler.Cliente = handler.NewClienteHandler(d.Service.Cliente)

		instance = d
	})
}
