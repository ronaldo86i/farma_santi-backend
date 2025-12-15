package port

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/johnfercher/maroto/v2/pkg/core"
)

type ReporteService interface {
	ReporteUsuariosPDF(ctx context.Context, filtros map[string]string) (core.Document, error)
	ReporteClientesPDF(ctx context.Context, filtros map[string]string) (core.Document, error)
	ReporteComprasPDF(ctx context.Context, filtros map[string]string) (core.Document, error)
	ReporteVentasPDF(ctx context.Context, filtros map[string]string) (core.Document, error)
	ReporteInventarioPDF(ctx context.Context, filtros map[string]string) (core.Document, error)
	ReporteLotesProductosPDF(ctx context.Context, filtros map[string]string) (core.Document, error)
	ReporteMovimientosPDF(ctx context.Context, filtros map[string]string) (core.Document, error)
	ReporteKardexProductoPDF(ctx context.Context, productoId *uuid.UUID) (core.Document, error)
	ReporteComprasDetallePDF(ctx context.Context, compraId *int) (core.Document, error)
}

type ReporteHandler interface {
	ReporteUsuariosPDF(c *fiber.Ctx) error
	ReporteClientesPDF(c *fiber.Ctx) error
	ReporteComprasPDF(c *fiber.Ctx) error
	ReporteVentasPDF(c *fiber.Ctx) error
	ReporteInventarioPDF(c *fiber.Ctx) error
	ReporteLotesProductosPDF(c *fiber.Ctx) error
	ReporteMovimientosPDF(c *fiber.Ctx) error
	ReporteKardexProductoPDF(c *fiber.Ctx) error
	ReporteComprasDetallePDF(c *fiber.Ctx) error
}
