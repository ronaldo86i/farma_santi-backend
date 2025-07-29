package port

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
)

type CompraRepository interface {
	RegistrarOrdenCompra(ctx context.Context, request *domain.CompraRequest) (*uint, error)
	ModificarOrdenCompra(ctx context.Context, id *int, request *domain.CompraRequest) error
	AnularOrdenCompra(ctx context.Context, id *int) error
	RegistrarCompra(ctx context.Context, id *int) error
	ObtenerListaCompras(ctx context.Context) (*[]domain.CompraInfo, error)
	ObtenerCompraById(ctx context.Context, id *int) (*domain.CompraDetail, error)
}

type CompraService interface {
	RegistrarOrdenCompra(ctx context.Context, request *domain.CompraRequest) (*uint, error)
	ModificarOrdenCompra(ctx context.Context, id *int, request *domain.CompraRequest) error
	AnularOrdenCompra(ctx context.Context, id *int) error
	RegistrarCompra(ctx context.Context, id *int) error
	ObtenerListaCompras(ctx context.Context) (*[]domain.CompraInfo, error)
	ObtenerCompraById(ctx context.Context, id *int) (*domain.CompraDetail, error)
}

type CompraHandler interface {
	RegistrarOrdenCompra(c *fiber.Ctx) error
	ModificarOrdenCompra(c *fiber.Ctx) error
	AnularOrdenCompra(c *fiber.Ctx) error
	RegistrarCompra(c *fiber.Ctx) error
	ObtenerListaCompras(c *fiber.Ctx) error
	ObtenerCompraById(c *fiber.Ctx) error
}
