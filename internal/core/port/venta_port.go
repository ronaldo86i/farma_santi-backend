package port

import (
	"context"
	"farma-santi_backend/internal/core/domain"

	"github.com/gofiber/fiber/v2"
)

type VentaRepository interface {
	ObtenerListaVentas(ctx context.Context, filtros map[string]string) (*[]domain.VentaInfo, error)
	RegistraVenta(ctx context.Context, request *domain.VentaRequest) (*int64, error)
	ObtenerVentaById(ctx context.Context, id *int) (*domain.VentaDetail, error)
	AnularVentaById(ctx context.Context, id *int) error
	FacturarVentaById(ctx context.Context, ventaId *int, req *domain.FacturaCompraVentaResponse) error
	ObtenerFacturaByVentaId(ctx context.Context, ventaId *int) (*domain.Factura, error)
}

type VentaService interface {
	ObtenerListaVentas(ctx context.Context, filtros map[string]string) (*[]domain.VentaInfo, error)
	RegistraVenta(ctx context.Context, request *domain.VentaRequest) (*int64, error)
	ObtenerVentaById(ctx context.Context, id *int) (*domain.VentaDetail, error)
	AnularVentaById(ctx context.Context, id *int) error
}

type VentaHandler interface {
	ObtenerListaVentas(c *fiber.Ctx) error
	RegistrarVenta(c *fiber.Ctx) error
	ObtenerVentaById(c *fiber.Ctx) error
	AnularVentaById(c *fiber.Ctx) error
	FacturarVentaById(c *fiber.Ctx) error
	ObtenerListaVentasShared(c *fiber.Ctx) error
	ObtenerVentaByIdShared(c *fiber.Ctx) error
}
