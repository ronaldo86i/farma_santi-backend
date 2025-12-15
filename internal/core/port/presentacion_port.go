package port

import (
	"context"
	"farma-santi_backend/internal/core/domain"

	"github.com/gofiber/fiber/v2"
)

type PresentacionRepository interface {
	ObtenerListaPresentaciones(ctx context.Context) (*[]domain.Presentacion, error)
}

type PresentacionService interface {
	ObtenerListaPresentaciones(ctx context.Context) (*[]domain.Presentacion, error)
}

type PresentacionHandler interface {
	ObtenerListaPresentaciones(c *fiber.Ctx) error
}
