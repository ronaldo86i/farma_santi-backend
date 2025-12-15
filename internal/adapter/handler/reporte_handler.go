package handler

import (
	"errors"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReporteHandler struct {
	reporteService port.ReporteService
}

func (r ReporteHandler) ReporteComprasDetallePDF(c *fiber.Ctx) error {
	compraId, err := c.ParamsInt("compraId", 0)
	if err != nil || compraId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' de la compra debe ser un número válido mayor a 0"))
	}
	doc, err := r.reporteService.ReporteComprasDetallePDF(c.UserContext(), &compraId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}

	c.Response().Header.Set("Content-Type", "application/pdf")
	c.Response().Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="reporte-compra-%d.pdf"`, compraId))
	c.Response().Header.Set("Content-Transfer-Encoding", "binary")

	return c.Send(doc.GetBytes())

}

func (r ReporteHandler) ReporteKardexProductoPDF(c *fiber.Ctx) error {
	productoIdParam := c.Params("productoId")
	productoId, err := uuid.Parse(productoIdParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(util.NewMessage("Formato de id no válido"))
	}
	doc, err := r.reporteService.ReporteKardexProductoPDF(c.UserContext(), &productoId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}

	c.Response().Header.Set("Content-Type", "application/pdf")
	c.Response().Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="reporte-kardex-%s.pdf"`, productoIdParam))
	c.Response().Header.Set("Content-Transfer-Encoding", "binary")

	return c.Send(doc.GetBytes())
}

func (r ReporteHandler) ReporteMovimientosPDF(c *fiber.Ctx) error {
	doc, err := r.reporteService.ReporteMovimientosPDF(c.UserContext(), c.Queries())
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}

	c.Response().Header.Set("Content-Type", "application/pdf")
	c.Response().Header.Set("Content-Disposition", "inline; filename=reporte-usuarios.pdf")
	c.Response().Header.Set("Content-Transfer-Encoding", "binary")

	return c.Send(doc.GetBytes())
}

func (r ReporteHandler) ReporteUsuariosPDF(c *fiber.Ctx) error {
	doc, err := r.reporteService.ReporteUsuariosPDF(c.UserContext(), c.Queries())
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}

	c.Response().Header.Set("Content-Type", "application/pdf")
	c.Response().Header.Set("Content-Disposition", "inline; filename=reporte-usuarios.pdf")
	c.Response().Header.Set("Content-Transfer-Encoding", "binary")

	return c.Send(doc.GetBytes())
}

func (r ReporteHandler) ReporteClientesPDF(c *fiber.Ctx) error {
	doc, err := r.reporteService.ReporteClientesPDF(c.UserContext(), c.Queries())
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}

	c.Response().Header.Set("Content-Type", "application/pdf")
	c.Response().Header.Set("Content-Disposition", "inline; filename=reporte-clientes.pdf")
	c.Response().Header.Set("Content-Transfer-Encoding", "binary")

	return c.Send(doc.GetBytes())
}

func (r ReporteHandler) ReporteComprasPDF(c *fiber.Ctx) error {
	doc, err := r.reporteService.ReporteComprasPDF(c.UserContext(), c.Queries())
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}

	c.Response().Header.Set("Content-Type", "application/pdf")
	c.Response().Header.Set("Content-Disposition", "inline; filename=reporte-compras.pdf")
	c.Response().Header.Set("Content-Transfer-Encoding", "binary")

	return c.Send(doc.GetBytes())
}

func (r ReporteHandler) ReporteVentasPDF(c *fiber.Ctx) error {
	doc, err := r.reporteService.ReporteVentasPDF(c.UserContext(), c.Queries())
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}

	c.Response().Header.Set("Content-Type", "application/pdf")
	c.Response().Header.Set("Content-Disposition", "inline; filename=reporte-ventas.pdf")
	c.Response().Header.Set("Content-Transfer-Encoding", "binary")

	return c.Send(doc.GetBytes())
}

func (r ReporteHandler) ReporteInventarioPDF(c *fiber.Ctx) error {
	doc, err := r.reporteService.ReporteInventarioPDF(c.UserContext(), c.Queries())
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}

	c.Response().Header.Set("Content-Type", "application/pdf")
	c.Response().Header.Set("Content-Disposition", "inline; filename=reporte-sinventario.pdf")
	c.Response().Header.Set("Content-Transfer-Encoding", "binary")

	return c.Send(doc.GetBytes())
}

func (r ReporteHandler) ReporteLotesProductosPDF(c *fiber.Ctx) error {
	doc, err := r.reporteService.ReporteLotesProductosPDF(c.UserContext(), c.Queries())
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
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
