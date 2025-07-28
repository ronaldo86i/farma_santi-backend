package handler

import (
	"errors"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/util"
	"log"
	"net/http"
	"strconv"

	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/port"
	"github.com/gofiber/fiber/v2"
)

type PrincipioActivoHandler struct {
	principioActivoService port.PrincipioActivoService
}

func (p PrincipioActivoHandler) RegistrarPrincipioActivo(c *fiber.Ctx) error {
	var req domain.PrincipioActivoRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("Petición inválida: datos incompletos o incorrectos"))
	}

	id, err := p.principioActivoService.RegistrarPrincipioActivo(c.Context(), &req)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}

	return c.Status(fiber.StatusCreated).JSON(util.NewMessageData(&domain.PrincipioActivoId{Id: *id}, "Principio activo registrado correctamente"))
}

func (p PrincipioActivoHandler) ModificarPrincipioActivo(c *fiber.Ctx) error {
	principioActivoId, err := c.ParamsInt("principioActivoId", 0)
	if err != nil || principioActivoId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del principio activo debe ser un número válido mayor a 0"))
	}

	var req domain.PrincipioActivoRequest
	err = c.BodyParser(&req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("Petición inválida: datos incompletos o incorrectos"))
	}

	err = p.principioActivoService.ModificarPrincipioActivo(c.Context(), &principioActivoId, &req)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}

	return c.JSON(util.NewMessage("Principio activo modificado correctamente"))
}

func (p PrincipioActivoHandler) ListarPrincipioActivo(c *fiber.Ctx) error {
	lista, err := p.principioActivoService.ListarPrincipioActivo(c.Context())
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}

	return c.JSON(lista)
}

func (p PrincipioActivoHandler) ObtenerPrincipioActivoById(c *fiber.Ctx) error {
	idParam := c.Params("principioActivoId")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id inválido"})
	}

	detalle, err := p.principioActivoService.ObtenerPrincipioActivoById(c.Context(), &id)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}
	return c.JSON(detalle)
}

func NewPrincipioActivoHandler(principioActivoService port.PrincipioActivoService) *PrincipioActivoHandler {
	return &PrincipioActivoHandler{principioActivoService: principioActivoService}
}

var _ port.PrincipioActivoHandler = (*PrincipioActivoHandler)(nil)
