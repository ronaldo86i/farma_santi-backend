package handler

import (
	"errors"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type StatHandler struct {
	statService port.StatService
}

func (s StatHandler) ObtenerEstadisticasDashboard(c *fiber.Ctx) error {
	stats, err := s.statService.ObtenerEstadisticasDashboard(c.UserContext())
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(stats)
}

func (s StatHandler) ObtenerTopProductosVendidos(c *fiber.Ctx) error {
	list, err := s.statService.ObtenerTopProductosVendidos(c.UserContext(), c.Queries())
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(list)
}

func NewStatHandler(statService port.StatService) *StatHandler {
	return &StatHandler{statService: statService}
}

var _ port.StatHandler = (*StatHandler)(nil)
