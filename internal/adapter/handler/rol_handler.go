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

type RolHandler struct {
	rolService port.RolService
}

func (r RolHandler) ModificarRol(c *fiber.Ctx) error {

	ctx := c.UserContext() // Usar el contexto

	var rolRequestUpdate domain.RolRequestUpdate
	if err := c.BodyParser(&rolRequestUpdate); err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("Petición inválida: datos incompletos o incorrectos"))
	}

	rolId, err := c.ParamsInt("rolId")
	if err != nil || rolId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del rol debe ser un número válido mayor a 0"))
	}
	err = r.rolService.ModificarRol(ctx, &rolId, &rolRequestUpdate)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}

	return c.Status(http.StatusAccepted).JSON(util.NewMessage("Rol actualizado correctamente"))
}

func (r RolHandler) RegistrarRol(c *fiber.Ctx) error {
	var rolRequest domain.RolRequest
	if err := c.BodyParser(&rolRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("Petición inválida: datos incompletos o incorrectos"))
	}

	ctx := c.UserContext() // Usar el contexto correcto

	err := r.rolService.RegistrarRol(ctx, &rolRequest)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(util.NewMessage("Rol registrado correctamente"))
}

func (r RolHandler) ModificarEstadoRol(c *fiber.Ctx) error {
	ctx := c.UserContext() // Usa el contexto de la petición
	rolId, err := c.ParamsInt("rolId")
	if err != nil || rolId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del rol debe ser un número válido mayor a 0"))
	}

	err = r.rolService.ModificarEstadoRol(ctx, &rolId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(util.NewMessage("Estado de rol actualizado correctamente"))
}

func (r RolHandler) ListarRoles(c *fiber.Ctx) error {
	ctx := c.UserContext() // Usa el contexto de la petición

	roles, err := r.rolService.ListarRoles(ctx)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(*roles)
}

func (r RolHandler) ObtenerRolById(c *fiber.Ctx) error {
	ctx := c.UserContext() // Usa el contexto de la petición

	rolId, err := c.ParamsInt("rolId")
	if err != nil || rolId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del rol debe ser un número válido mayor a 0"))
	}

	rol, err := r.rolService.ObtenerRolById(ctx, &rolId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage("Error interno del servidor"))
	}

	return c.Status(http.StatusOK).JSON(&rol)
}

func NewRolHandler(rolService port.RolService) *RolHandler {
	return &RolHandler{rolService}
}

var _ port.RolHandler = (*RolHandler)(nil)
