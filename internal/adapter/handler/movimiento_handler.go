package handler

import (
	"errors"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
	"log"

	"github.com/gofiber/fiber/v2"
)

type MovimientoHandler struct {
	movimientoService port.MovimientoService
}

func (m MovimientoHandler) ObtenerMovimientosKardex(c *fiber.Ctx) error {
	list, err := m.movimientoService.ObtenerMovimientosKardex(c.UserContext(), c.Queries())
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}
	return c.JSON(&list)
}

func (m MovimientoHandler) ObtenerListaMovimientos(c *fiber.Ctx) error {
	list, err := m.movimientoService.ObtenerListaMovimientos(c.UserContext(), c.Queries())
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}
	return c.JSON(&list)
}

func NewMovimientoHandler(movimientoService port.MovimientoService) *MovimientoHandler {
	return &MovimientoHandler{movimientoService: movimientoService}
}

var _ port.MovimientoHandler = (*MovimientoHandler)(nil)
