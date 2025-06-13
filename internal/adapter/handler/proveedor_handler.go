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

type ProveedorHandler struct {
	proveedorService port.ProveedorService
}

func (p ProveedorHandler) HabilitarProveedor(c *fiber.Ctx) error {
	ctx := c.UserContext()
	proveedorId, err := c.ParamsInt("proveedorId", 0)
	if err != nil || proveedorId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del proveedor debe ser un número válido mayor a 0"))
	}
	err = p.proveedorService.HabilitarProveedor(ctx, &proveedorId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}
	return c.Status(http.StatusOK).JSON(util.NewMessage("Proveedor actualizado correctamente"))
}

func (p ProveedorHandler) DeshabilitarProveedor(c *fiber.Ctx) error {
	ctx := c.UserContext()
	proveedorId, err := c.ParamsInt("proveedorId", 0)
	if err != nil || proveedorId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del proveedor debe ser un número válido mayor a 0"))
	}
	err = p.proveedorService.DeshabilitarProveedor(ctx, &proveedorId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}
	return c.Status(http.StatusOK).JSON(util.NewMessage("Proveedor actualizado correctamente"))
}

func (p ProveedorHandler) RegistrarProveedor(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var proveedor domain.ProveedorRequest
	if err := c.BodyParser(&proveedor); err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("Petición inválida: datos incompletos o incorrectos"))
	}
	err := p.proveedorService.RegistrarProveedor(ctx, &proveedor)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}
	return c.Status(http.StatusOK).JSON(util.NewMessage("Proveedor registrado correctamente"))
}

func (p ProveedorHandler) ObtenerProveedorById(c *fiber.Ctx) error {
	ctx := c.UserContext()
	proveedorId, err := c.ParamsInt("proveedorId", 0)
	if err != nil || proveedorId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del proveedor debe ser un número válido mayor a 0"))
	}
	proveedor, err := p.proveedorService.ObtenerProveedorById(ctx, &proveedorId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}

	return c.Status(http.StatusOK).JSON(&proveedor)
}

func (p ProveedorHandler) ListarProveedores(c *fiber.Ctx) error {
	ctx := c.UserContext()
	proveedores, err := p.proveedorService.ListarProveedores(ctx)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}
	return c.Status(http.StatusOK).JSON(*proveedores)
}

func (p ProveedorHandler) ModificarProveedor(c *fiber.Ctx) error {
	ctx := c.UserContext()
	proveedorId, err := c.ParamsInt("proveedorId", 0)
	if err != nil || proveedorId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del proveedor debe ser un número válido mayor a 0"))
	}

	var proveedor domain.ProveedorRequest
	if err := c.BodyParser(&proveedor); err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("Petición inválida: datos incompletos o incorrectos"))
	}

	err = p.proveedorService.ModificarProveedor(ctx, &proveedorId, &proveedor)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}

	return c.Status(http.StatusOK).JSON(util.NewMessage("Proveedor actualizado correctamente"))
}

func NewProveedorHandler(proveedorService port.ProveedorService) *ProveedorHandler {
	return &ProveedorHandler{proveedorService}
}

var _ port.ProveedorHandler = (*ProveedorHandler)(nil)
