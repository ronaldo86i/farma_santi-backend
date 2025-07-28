package port

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
)

type MovimientoRepository interface {
	ObtenerListaMovimientos(ctx context.Context) (*[]domain.MovimientoInfo, error)
}

type MovimientoService interface {
	ObtenerListaMovimientos(ctx context.Context) (*[]domain.MovimientoInfo, error)
}

type MovimientoHandler interface {
	ObtenerListaMovimientos(c *fiber.Ctx) error
}
