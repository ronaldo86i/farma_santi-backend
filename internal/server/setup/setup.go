package setup

import (
	"farma-santi_backend/internal/adapter/handler"
	"farma-santi_backend/internal/adapter/repository"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/service"
	"farma-santi_backend/internal/postgresql"
	"log"
	"sync"

	"github.com/joho/godotenv"
)

type Repository struct {
	Categoria       port.CategoriaRepository
	Cliente         port.ClienteRepository
	Compra          port.CompraRepository
	Laboratorio     port.LaboratorioRepository
	LoteProducto    port.LoteProductoRepository
	PrincipioActivo port.PrincipioActivoRepository
	Producto        port.ProductoRepository
	Proveedor       port.ProveedorRepository
	Rol             port.RolRepository
	Usuario         port.UsuarioRepository
	Venta           port.VentaRepository
	Movimiento      port.MovimientoRepository
	Presentacion    port.PresentacionRepository
	Stat            port.StatRepository
}

type Service struct {
	Auth            port.AuthService
	Categoria       port.CategoriaService
	Cliente         port.ClienteService
	Compra          port.CompraService
	Laboratorio     port.LaboratorioService
	LoteProducto    port.LoteProductoService
	PrincipioActivo port.PrincipioActivoService
	Producto        port.ProductoService
	Proveedor       port.ProveedorService
	Rol             port.RolService
	Usuario         port.UsuarioService
	Venta           port.VentaService
	Movimiento      port.MovimientoService
	Reporte         port.ReporteService
	Presentacion    port.PresentacionService
	Stat            port.StatService
	Backup          port.BackupService
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
	Venta           port.VentaHandler
	Movimiento      port.MovimientoHandler
	Reporte         port.ReporteHandler
	Presentacion    port.PresentacionHandler
	Stat            port.StatHandler
	Backup          port.BackupHandler
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
		repositories := &d.Repository
		services := &d.Service
		handlers := &d.Handler

		// Repositories
		repositories.Usuario = repository.NewUsuarioRepository(pool)
		repositories.Rol = repository.NewRolRepository(pool)
		repositories.Categoria = repository.NewCategoriaRepository(pool)
		repositories.Proveedor = repository.NewProveedorRepository(pool)
		repositories.Laboratorio = repository.NewLaboratorioRepository(pool)
		repositories.Producto = repository.NewProductoRepository(pool)
		repositories.LoteProducto = repository.NewLoteProductoRepository(pool)
		repositories.PrincipioActivo = repository.NewPrincipioActivoRepository(pool)
		repositories.Compra = repository.NewCompraRepository(pool)
		repositories.Cliente = repository.NewClienteRepository(pool)
		repositories.Venta = repository.NewVentaRepository(pool)
		repositories.Movimiento = repository.NewMovimientoRepository(pool)
		repositories.Presentacion = repository.NewPresentacionRepository(pool)
		repositories.Stat = repository.NewStatRepository(pool)
		// Services
		services.Auth = service.NewAuthService(repositories.Usuario)
		services.Usuario = service.NewUsuarioService(repositories.Usuario)
		services.Rol = service.NewRolService(repositories.Rol)
		services.Categoria = service.NewCategoriaService(repositories.Categoria)
		services.Proveedor = service.NewProveedorService(repositories.Proveedor)
		services.Laboratorio = service.NewLaboratorioService(repositories.Laboratorio)
		services.Producto = service.NewProductoService(repositories.Producto)
		services.LoteProducto = service.NewLoteProductoService(repositories.LoteProducto)
		services.PrincipioActivo = service.NewPrincipioActivoService(repositories.PrincipioActivo)
		services.Compra = service.NewCompraService(repositories.Compra)
		services.Cliente = service.NewClienteService(repositories.Cliente)
		services.Venta = service.NewVentaService(repositories.Venta)
		services.Movimiento = service.NewMovimientoService(repositories.Movimiento)
		services.Reporte = service.NewReporteService(repositories.Usuario, repositories.Cliente, repositories.LoteProducto, repositories.Producto, repositories.Compra, repositories.Venta, repositories.Movimiento)
		services.Presentacion = service.NewPresentacionService(repositories.Presentacion)
		services.Stat = service.NewStatService(repositories.Stat)
		services.Backup = service.NewBackupService()
		// Handlers
		handlers.Auth = handler.NewAuthHandler(services.Auth)
		handlers.Usuario = handler.NewUsuarioHandler(services.Usuario)
		handlers.Rol = handler.NewRolHandler(services.Rol)
		handlers.Categoria = handler.NewCategoriaHandler(services.Categoria)
		handlers.Proveedor = handler.NewProveedorHandler(services.Proveedor)
		handlers.Laboratorio = handler.NewLaboratorioHandler(services.Laboratorio)
		handlers.Producto = handler.NewProductoHandler(services.Producto)
		handlers.LoteProducto = handler.NewLoteProductoHandler(services.LoteProducto)
		handlers.PrincipioActivo = handler.NewPrincipioActivoHandler(services.PrincipioActivo)
		handlers.Compra = handler.NewCompraHandler(services.Compra)
		handlers.Cliente = handler.NewClienteHandler(services.Cliente)
		handlers.Venta = handler.NewVentaHandler(services.Venta)
		handlers.Movimiento = handler.NewMovimientoHandler(services.Movimiento)
		handlers.Reporte = handler.NewReporteHandler(services.Reporte)
		handlers.Presentacion = handler.NewPresentacionHandler(services.Presentacion)
		handlers.Stat = handler.NewStatHandler(services.Stat)
		handlers.Backup = handler.NewBackupHandler(services.Backup)

		instance = d
	})
}
