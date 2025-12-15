package port

import (
	"context"
	"farma-santi_backend/internal/core/domain"

	"github.com/gofiber/fiber/v2"
)

type StatRepository interface {
	ObtenerTopProductosVendidos(ctx context.Context, filtros map[string]string) (*[]domain.ProductoStat, error)
	ObtenerEstadisticasDashboard(ctx context.Context) (*domain.DashboardStats, error)
}

type StatService interface {
	ObtenerTopProductosVendidos(ctx context.Context, filtros map[string]string) (*[]domain.ProductoStat, error)
	ObtenerEstadisticasDashboard(ctx context.Context) (*domain.DashboardStats, error)
}

type StatHandler interface {
	ObtenerTopProductosVendidos(c *fiber.Ctx) error
	ObtenerEstadisticasDashboard(c *fiber.Ctx) error
}
