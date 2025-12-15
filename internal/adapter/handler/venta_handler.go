package handler

import (
	"context"
	"errors"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/service"
	"farma-santi_backend/internal/core/util"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type VentaHandler struct {
	ventaService port.VentaService
}

func (v VentaHandler) ObtenerVentaByIdShared(c *fiber.Ctx) error {
	ventaId, err := c.ParamsInt("ventaId")
	if err != nil || ventaId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' de venta debe ser un número válido mayor a 0"))
	}

	val := c.UserContext().Value(util.ContextUserIdKey)
	userId, ok := val.(string)
	if !ok {
		return datatype.NewBadRequestError("UID de usuario inválido o no encontrado en el contexto")
	}
	fb := service.GetFirebaseClient()
	token, err := fb.AuthClient.GetUser(context.Background(), userId)
	if err != nil {
		log.Println("Error al verificar uid de usuario", err)
		return datatype.NewInternalServerErrorGeneric()
	}
	email := token.Email

	venta, err := v.ventaService.ObtenerVentaById(c.UserContext(), &ventaId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	if venta.Cliente.Email != email {
		return c.Status(http.StatusNotFound).JSON("Compra no encontrada")
	}
	return c.JSON(&venta)
}

func (v VentaHandler) ObtenerListaVentasShared(c *fiber.Ctx) error {
	val := c.UserContext().Value(util.ContextUserIdKey)
	userId, ok := val.(string)
	if !ok {
		return datatype.NewBadRequestError("UID de usuario inválido o no encontrado en el contexto")
	}
	fb := service.GetFirebaseClient()
	token, err := fb.AuthClient.GetUser(context.Background(), userId)
	if err != nil {
		log.Println("Error al verificar uid de usuario", err)
		return datatype.NewInternalServerErrorGeneric()
	}
	email := token.Email
	filtros := map[string]string{
		"email":  email,
		"estado": "Realizada",
	}
	list, err := v.ventaService.ObtenerListaVentas(c.UserContext(), filtros)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(&list)
}

func (v VentaHandler) ObtenerListaVentas(c *fiber.Ctx) error {
	list, err := v.ventaService.ObtenerListaVentas(c.UserContext(), c.Queries())
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(&list)
}

func (v VentaHandler) RegistrarVenta(c *fiber.Ctx) error {
	var venta *domain.VentaRequest
	err := c.BodyParser(&venta)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("Petición inválida: datos incompletos o incorrectos"))
	}
	ventaId, err := v.ventaService.RegistraVenta(c.UserContext(), venta)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		var errorDataResponse *datatype.ErrorDataResponse[domain.ProductoId]

		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		} else if errors.As(err, &errorDataResponse) {
			return c.Status(errorDataResponse.Code).JSON(&errorDataResponse)
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.Status(http.StatusCreated).JSON(util.NewMessageData(domain.VentaResponse{VentaId: *ventaId}, "Venta registrada correctamente"))
}

func (v VentaHandler) ObtenerVentaById(c *fiber.Ctx) error {
	ventaId, err := c.ParamsInt("ventaId")
	if err != nil || ventaId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' de venta debe ser un número válido mayor a 0"))
	}
	venta, err := v.ventaService.ObtenerVentaById(c.UserContext(), &ventaId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(&venta)
}

func (v VentaHandler) AnularVentaById(c *fiber.Ctx) error {
	ventaId, err := c.ParamsInt("ventaId")
	if err != nil || ventaId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' de venta debe ser un número válido mayor a 0"))
	}
	err = v.ventaService.AnularVentaById(c.UserContext(), &ventaId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.Status(http.StatusOK).JSON(util.NewMessage("Venta anulada correctamente"))
}

func (v VentaHandler) FacturarVentaById(c *fiber.Ctx) error {
	ventaId, err := c.ParamsInt("ventaId")
	if err != nil || ventaId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' de venta debe ser un número válido mayor a 0"))
	}
	err = v.ventaService.AnularVentaById(c.UserContext(), &ventaId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.Status(http.StatusOK).JSON(util.NewMessage("Venta anulada correctamente"))
}

func NewVentaHandler(ventaService port.VentaService) *VentaHandler {
	return &VentaHandler{ventaService: ventaService}
}

var _ port.VentaHandler = (*VentaHandler)(nil)
