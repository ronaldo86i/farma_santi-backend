package port

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/johnfercher/maroto/v2/pkg/core"
)

type ReporteService interface {
	ReporteUsuariosPDF(ctx context.Context) (core.Document, error)
	ReporteClientesPDF(ctx context.Context) (core.Document, error)
	ReporteComprasPDF(ctx context.Context) (core.Document, error)
	ReporteVentasPDF(ctx context.Context) (core.Document, error)
	ReporteInventarioPDF(ctx context.Context) (core.Document, error)
	ReporteLotesProductosPDF(ctx context.Context) (core.Document, error)
}

type ReporteHandler interface {
	ReporteUsuariosPDF(c *fiber.Ctx) error
	ReporteClientesPDF(c *fiber.Ctx) error
	ReporteComprasPDF(c *fiber.Ctx) error
	ReporteVentasPDF(c *fiber.Ctx) error
	ReporteInventarioPDF(c *fiber.Ctx) error
	ReporteLotesProductosPDF(c *fiber.Ctx) error
}
