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

type ClienteHandler struct {
	clienteService port.ClienteService
}

func (c2 ClienteHandler) ObtenerListaClientes(c *fiber.Ctx) error {
	list, err := c2.clienteService.ObtenerListaClientes(c.UserContext())
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

func (c2 ClienteHandler) ObtenerClienteById(c *fiber.Ctx) error {
	clienteId, err := c.ParamsInt("clienteId", 0)
	if err != nil || clienteId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del cliente debe ser un número válido mayor a 0"))
	}
	cliente, err := c2.clienteService.ObtenerClienteById(c.UserContext(), &clienteId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(&cliente)
}

func (c2 ClienteHandler) RegistrarCliente(c *fiber.Ctx) error {
	var request domain.ClienteRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("Petición inválida: datos incompletos o incorrectos"))
	}
	clienteId, err := c2.clienteService.RegistrarCliente(c.UserContext(), &request)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.Status(http.StatusOK).JSON(util.NewMessageData(domain.ClienteId{Id: *clienteId}, "Cliente registrado correctamente"))
}

func (c2 ClienteHandler) ModificarClienteById(c *fiber.Ctx) error {
	clienteId, err := c.ParamsInt("clienteId", 0)
	if err != nil || clienteId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del cliente debe ser un número válido mayor a 0"))
	}
	var request domain.ClienteRequest
	err = c.BodyParser(&request)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("Petición inválida: datos incompletos o incorrectos"))
	}
	err = c2.clienteService.ModificarClienteById(c.UserContext(), &clienteId, &request)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.Status(http.StatusOK).JSON(util.NewMessage("Cliente modificado correctamente"))
}

func (c2 ClienteHandler) HabilitarCliente(c *fiber.Ctx) error {
	clienteId, err := c.ParamsInt("clienteId", 0)
	if err != nil || clienteId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del cliente debe ser un número válido mayor a 0"))
	}
	err = c2.clienteService.HabilitarCliente(c.UserContext(), &clienteId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.Status(http.StatusOK).JSON(util.NewMessage("Cliente actualizado correctamente"))
}

func (c2 ClienteHandler) DeshabilitarCliente(c *fiber.Ctx) error {
	clienteId, err := c.ParamsInt("clienteId", 0)
	if err != nil || clienteId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del cliente debe ser un número válido mayor a 0"))
	}
	err = c2.clienteService.DeshabilitarCliente(c.UserContext(), &clienteId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.Status(http.StatusOK).JSON(util.NewMessage("Cliente actualizado correctamente"))
}

func NewClienteHandler(clienteService port.ClienteService) ClienteHandler {
	return ClienteHandler{clienteService: clienteService}
}

var _ port.ClienteHandler = (*ClienteHandler)(nil)
