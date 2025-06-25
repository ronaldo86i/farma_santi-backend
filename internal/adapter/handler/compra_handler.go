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

type CompraHandler struct {
	compraService port.CompraService
}

func (c2 CompraHandler) ObtenerCompraById(c *fiber.Ctx) error {
	compraId, err := c.ParamsInt("compraId", 0)
	if err != nil || compraId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' de la compra debe ser un número válido mayor a 0"))
	}
	compra, err := c2.compraService.ObtenerCompraById(c.UserContext(), &compraId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.Status(http.StatusOK).JSON(&compra)
}

func (c2 CompraHandler) RegistrarOrdenCompra(c *fiber.Ctx) error {
	var request domain.CompraRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("Petición inválida: datos incompletos o incorrectos"))
	}
	err = c2.compraService.RegistrarOrdenCompra(c.UserContext(), &request)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(util.NewMessage("Orden de compra registrada correctamente"))
}

func (c2 CompraHandler) ModificarOrdenCompra(c *fiber.Ctx) error {
	var request domain.CompraRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("Petición inválida: datos incompletos o incorrectos"))
	}
	compraId, err := c.ParamsInt("compraId", 0)
	if err != nil || compraId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' de la compra debe ser un número válido mayor a 0"))
	}
	err = c2.compraService.ModificarOrdenCompra(c.UserContext(), &compraId, &request)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(util.NewMessage("Orden de compra modificada correctamente"))
}

func (c2 CompraHandler) AnularOrdenCompra(c *fiber.Ctx) error {
	compraId, err := c.ParamsInt("compraId", 0)
	if err != nil || compraId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' de la compra debe ser un número válido mayor a 0"))
	}
	err = c2.compraService.AnularOrdenCompra(c.UserContext(), &compraId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(util.NewMessage("Orden de compra anulada correctamente"))
}

func (c2 CompraHandler) RegistrarCompra(c *fiber.Ctx) error {
	compraId, err := c.ParamsInt("compraId", 0)
	if err != nil || compraId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' de la compra debe ser un número válido mayor a 0"))
	}
	err = c2.compraService.RegistrarCompra(c.UserContext(), &compraId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(util.NewMessage("Compra registrada correctamente"))
}

func (c2 CompraHandler) ObtenerListaCompras(c *fiber.Ctx) error {
	lista, err := c2.compraService.ObtenerListaCompras(c.UserContext())
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(&lista)
}

func NewCompraHandler(compraService port.CompraService) *CompraHandler {
	return &CompraHandler{compraService: compraService}
}

var _ port.CompraHandler = (*CompraHandler)(nil)
