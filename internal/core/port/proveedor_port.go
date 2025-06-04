package port

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
)

type ProveedorRepository interface {
	RegistrarProveedor(ctx context.Context, request *domain.ProveedorRequest) error
	ObtenerProveedorById(ctx context.Context, id *int) (*domain.ProveedorDetail, error)
	ListarProveedores(ctx context.Context) (*[]domain.ProveedorInfo, error)
	ModificarProveedor(ctx context.Context, id *int, request *domain.ProveedorRequest) error
	HabilitarProveedor(ctx context.Context, id *int) error
	DeshabilitarProveedor(ctx context.Context, id *int) error
}

type ProveedorService interface {
	RegistrarProveedor(ctx context.Context, request *domain.ProveedorRequest) error
	ObtenerProveedorById(ctx context.Context, id *int) (*domain.ProveedorDetail, error)
	ListarProveedores(ctx context.Context) (*[]domain.ProveedorInfo, error)
	ModificarProveedor(ctx context.Context, id *int, request *domain.ProveedorRequest) error
	HabilitarProveedor(ctx context.Context, id *int) error
	DeshabilitarProveedor(ctx context.Context, id *int) error
}

type ProveedorHandler interface {
	RegistrarProveedor(c *fiber.Ctx) error
	ObtenerProveedorById(c *fiber.Ctx) error
	ListarProveedores(c *fiber.Ctx) error
	ModificarProveedor(c *fiber.Ctx) error
	HabilitarProveedor(c *fiber.Ctx) error
	DeshabilitarProveedor(c *fiber.Ctx) error
}
