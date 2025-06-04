package handler

import (
	"errors"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
)

type LaboratorioHandler struct {
	laboratorioService port.LaboratorioService
}

func (l LaboratorioHandler) ListarLaboratoriosDisponibles(c *fiber.Ctx) error {
	list, err := l.laboratorioService.ListarLaboratoriosDisponibles(c.UserContext())
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerError()
	}
	return c.JSON(&list)
}

func (l LaboratorioHandler) ListarLaboratorios(c *fiber.Ctx) error {
	list, err := l.laboratorioService.ListarLaboratorios(c.UserContext())
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerError()
	}
	return c.JSON(&list)
}

func (l LaboratorioHandler) ObtenerLaboratorioById(c *fiber.Ctx) error {
	ctx := c.UserContext()
	laboratorioId, err := c.ParamsInt("laboratorioId", 0)
	if err != nil || laboratorioId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del laboratorio debe ser un número válido mayor a 0"))
	}
	lab, err := l.laboratorioService.ObtenerLaboratorioById(ctx, &laboratorioId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerError()
	}
	return c.JSON(&lab)
}

func (l LaboratorioHandler) RegistrarLaboratorio(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var laboratorioRequest domain.LaboratorioRequest
	err := c.BodyParser(&laboratorioRequest)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("Petición inválida: datos incompletos o incorrectos"))
	}
	err = l.laboratorioService.RegistrarLaboratorio(ctx, &laboratorioRequest)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.Status(http.StatusOK).JSON(util.NewMessage("Laboratorio registrado correctamente"))
}

func (l LaboratorioHandler) ModificarLaboratorio(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var laboratorioRequest domain.LaboratorioRequest
	err := c.BodyParser(&laboratorioRequest)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("Petición inválida: datos incompletos o incorrectos"))
	}
	laboratorioId, err := c.ParamsInt("laboratorioId", 0)
	if err != nil || laboratorioId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del laboratorio debe ser un número válido mayor a 0"))
	}
	err = l.laboratorioService.ModificarLaboratorio(ctx, &laboratorioId, &laboratorioRequest)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.Status(http.StatusOK).JSON(util.NewMessage("Laboratorio modificado correctamente"))
}

func (l LaboratorioHandler) HabilitarLaboratorio(c *fiber.Ctx) error {
	ctx := c.UserContext()
	laboratorioId, err := c.ParamsInt("laboratorioId", 0)
	if err != nil || laboratorioId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del laboratorio debe ser un número válido mayor a 0"))
	}
	err = l.laboratorioService.HabilitarLaboratorio(ctx, &laboratorioId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerError()
	}
	return c.JSON(util.NewMessage("Laboratorio actualizado correctamente"))
}

func (l LaboratorioHandler) DeshabilitarLaboratorio(c *fiber.Ctx) error {
	ctx := c.UserContext()
	laboratorioId, err := c.ParamsInt("laboratorioId", 0)
	if err != nil || laboratorioId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del laboratorio debe ser un número válido mayor a 0"))
	}
	err = l.laboratorioService.DeshabilitarLaboratorio(ctx, &laboratorioId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerError()
	}
	return c.JSON(util.NewMessage("Laboratorio actualizado correctamente"))
}

func NewLaboratorioHandler(laboratorioService port.LaboratorioService) *LaboratorioHandler {
	return &LaboratorioHandler{laboratorioService: laboratorioService}
}

var _ port.LaboratorioHandler = (*LaboratorioHandler)(nil)
