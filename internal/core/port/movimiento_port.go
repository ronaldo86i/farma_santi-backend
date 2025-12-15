package port

import (
	"context"
	"farma-santi_backend/internal/core/domain"

	"github.com/gofiber/fiber/v2"
)

type MovimientoRepository interface {
	ObtenerListaMovimientos(ctx context.Context, filtros map[string]string) (*[]domain.MovimientoInfo, error)
	ObtenerMovimientosKardex(ctx context.Context, filtros map[string]string) (*[]domain.MovimientoKardex, error)
}

type MovimientoService interface {
	ObtenerListaMovimientos(ctx context.Context, filtros map[string]string) (*[]domain.MovimientoInfo, error)
	ObtenerMovimientosKardex(ctx context.Context, filtros map[string]string) (*[]domain.MovimientoKardex, error)
}

type MovimientoHandler interface {
	ObtenerListaMovimientos(c *fiber.Ctx) error
	ObtenerMovimientosKardex(c *fiber.Ctx) error
}
