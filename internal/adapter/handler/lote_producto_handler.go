package handler

import (
	"errors"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type LoteProductoHandler struct {
	loteProductoService port.LoteProductoService
}

func (l LoteProductoHandler) ListarLotesProductosByProductoId(c *fiber.Ctx) error {
	// Obtener id del producto del parámetro
	productoIdParam := c.Params("productoId")
	productoId, err := uuid.Parse(productoIdParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(util.NewMessage("Formato de id no válido"))
	}

	list, err := l.loteProductoService.ListarLotesProductosByProductoId(c.UserContext(), &productoId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}
	return c.Status(http.StatusOK).JSON(&list)
}

func (l LoteProductoHandler) ModificarLoteProducto(c *fiber.Ctx) error {
	loteProductoId, err := c.ParamsInt("loteProductoId", 0)
	if err != nil || loteProductoId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del proveedor debe ser un número válido mayor a 0"))
	}

	var loteProducto domain.LoteProductoRequest
	if err := c.BodyParser(&loteProducto); err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("Petición inválida: datos incompletos o incorrectos"))
	}
	err = l.loteProductoService.ModificarLoteProducto(c.UserContext(), &loteProductoId, &loteProducto)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}
	return c.Status(http.StatusCreated).JSON(util.NewMessage("Lote de producto modificado correctamente"))
}

func (l LoteProductoHandler) ListarLotesProductos(c *fiber.Ctx) error {
	list, err := l.loteProductoService.ListarLotesProductos(c.UserContext())
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}
	return c.Status(http.StatusOK).JSON(&list)
}

func (l LoteProductoHandler) RegistrarLoteProducto(c *fiber.Ctx) error {
	var loteProducto domain.LoteProductoRequest
	if err := c.BodyParser(&loteProducto); err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("Petición inválida: datos incompletos o incorrectos"))
	}
	err := l.loteProductoService.RegistrarLoteProducto(c.UserContext(), &loteProducto)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}
	return c.Status(http.StatusCreated).JSON(util.NewMessage("Lote de producto registrado correctamente"))
}

func (l LoteProductoHandler) ObtenerLoteProductoById(c *fiber.Ctx) error {
	loteProductoId, err := c.ParamsInt("loteProductoId", 0)
	if err != nil || loteProductoId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del proveedor debe ser un número válido mayor a 0"))
	}

	loteProducto, err := l.loteProductoService.ObtenerLoteProductoById(c.UserContext(), &loteProductoId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}

	return c.Status(http.StatusOK).JSON(&loteProducto)
}

func NewLoteProductoHandler(loteProductoService port.LoteProductoService) *LoteProductoHandler {
	return &LoteProductoHandler{loteProductoService}
}

var _ port.LoteProductoHandler = (*LoteProductoHandler)(nil)
