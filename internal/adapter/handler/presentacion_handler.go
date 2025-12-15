package handler

import (
	"errors"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
	"log"

	"github.com/gofiber/fiber/v2"
)

type PresentacionHandler struct {
	presentacionService port.PresentacionService
}

func (p PresentacionHandler) ObtenerListaPresentaciones(c *fiber.Ctx) error {
	list, err := p.presentacionService.ObtenerListaPresentaciones(c.UserContext())
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

func NewPresentacionHandler(presentacionService port.PresentacionService) *PresentacionHandler {
	return &PresentacionHandler{presentacionService: presentacionService}
}

var _ port.PresentacionHandler = (*PresentacionHandler)(nil)
