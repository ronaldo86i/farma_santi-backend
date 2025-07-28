package handler

import (
	"farma-santi_backend/internal/core/port"
	"github.com/gofiber/fiber/v2"
)

type ReporteHandler struct {
	reporteService port.ReporteService
}

func (r ReporteHandler) ReporteUsuariosPDF(c *fiber.Ctx) error {
	doc, err := r.reporteService.ReporteUsuariosPDF(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error generando el reporte",
		})
	}

	c.Response().Header.Set("Content-Type", "application/pdf")
	c.Response().Header.Set("Content-Disposition", "inline; filename=reporte-usuarios.pdf")
	c.Response().Header.Set("Content-Transfer-Encoding", "binary")

	return c.Send(doc.GetBytes())
}

func (r ReporteHandler) ReporteClientesPDF(c *fiber.Ctx) error {
	doc, err := r.reporteService.ReporteClientesPDF(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error generando el reporte",
		})
	}

	c.Response().Header.Set("Content-Type", "application/pdf")
	c.Response().Header.Set("Content-Disposition", "inline; filename=reporte-clientes.pdf")
	c.Response().Header.Set("Content-Transfer-Encoding", "binary")

	return c.Send(doc.GetBytes())
}

func (r ReporteHandler) ReporteComprasPDF(c *fiber.Ctx) error {
	doc, err := r.reporteService.ReporteComprasPDF(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error generando el reporte",
		})
	}

	c.Response().Header.Set("Content-Type", "application/pdf")
	c.Response().Header.Set("Content-Disposition", "inline; filename=reporte-compras.pdf")
	c.Response().Header.Set("Content-Transfer-Encoding", "binary")

	return c.Send(doc.GetBytes())
}

func (r ReporteHandler) ReporteVentasPDF(c *fiber.Ctx) error {
	doc, err := r.reporteService.ReporteVentasPDF(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error generando el reporte",
		})
	}

	c.Response().Header.Set("Content-Type", "application/pdf")
	c.Response().Header.Set("Content-Disposition", "inline; filename=reporte-ventas.pdf")
	c.Response().Header.Set("Content-Transfer-Encoding", "binary")

	return c.Send(doc.GetBytes())
}

func (r ReporteHandler) ReporteInventarioPDF(c *fiber.Ctx) error {
	doc, err := r.reporteService.ReporteInventarioPDF(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error generando el reporte",
		})
	}

	c.Response().Header.Set("Content-Type", "application/pdf")
	c.Response().Header.Set("Content-Disposition", "inline; filename=reporte-sinventario.pdf")
	c.Response().Header.Set("Content-Transfer-Encoding", "binary")

	return c.Send(doc.GetBytes())
}

func (r ReporteHandler) ReporteLotesProductosPDF(c *fiber.Ctx) error {
	doc, err := r.reporteService.ReporteLotesProductosPDF(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error generando el reporte",
		})
	}

	c.Response().Header.Set("Content-Type", "application/pdf")
	c.Response().Header.Set("Content-Disposition", "inline; filename=reporte-lotes-productos.pdf")
	c.Response().Header.Set("Content-Transfer-Encoding", "binary")

	return c.Send(doc.GetBytes())
}

func NewReporteHandler(reporteService port.ReporteService) *ReporteHandler {
	return &ReporteHandler{reporteService: reporteService}
}

var _ port.ReporteHandler = (*ReporteHandler)(nil)
