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

type CategoriaHandler struct {
	categoriaService port.CategoriaService
}

func (ch CategoriaHandler) ObtenerCategoriaById(c *fiber.Ctx) error {
	ctx := c.Context()
	categoriaId, err := c.ParamsInt("categoriaId", 0)
	if err != nil || categoriaId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del categoria debe ser un número válido mayor a 0"))
	}
	categoria, err := ch.categoriaService.ObtenerCategoriaById(ctx, &categoriaId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(&categoria)
}

func (ch CategoriaHandler) ListarCategorias(c *fiber.Ctx) error {
	ctx := c.Context()
	categorias, err := ch.categoriaService.ListarCategorias(ctx)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(&categorias)
}

func (ch CategoriaHandler) ModificarEstadoCategoria(c *fiber.Ctx) error {
	ctx := c.Context()
	categoriaId, err := c.ParamsInt("categoriaId", 0)
	if err != nil || categoriaId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del categoria debe ser un número válido mayor a 0"))
	}
	err = ch.categoriaService.ModificarEstadoCategoria(ctx, &categoriaId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(util.NewMessage("Categoría actualizada correctamente"))
}

func (ch CategoriaHandler) ModificarCategoria(c *fiber.Ctx) error {
	ctx := c.Context()
	categoriaId, err := c.ParamsInt("categoriaId", 0)
	if err != nil || categoriaId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del categoria debe ser un número válido mayor a 0"))
	}
	var categoriaRequest domain.CategoriaRequest
	if err := c.BodyParser(&categoriaRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("Petición inválida: datos incompletos o incorrectos"))
	}

	err = ch.categoriaService.ModificarCategoria(ctx, &categoriaId, &categoriaRequest)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(util.NewMessage("Categoría modificada correctamente"))
}

func (ch CategoriaHandler) RegistrarCategoria(c *fiber.Ctx) error {
	var categoriaRequest domain.CategoriaRequest
	if err := c.BodyParser(&categoriaRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("Petición inválida: datos incompletos o incorrectos"))
	}
	ctx := c.Context()
	err := ch.categoriaService.RegistrarCategoria(ctx, &categoriaRequest)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(util.NewMessage("Categoría registrado correctamente"))
}

func NewCategoriaHandler(categoriaService port.CategoriaService) *CategoriaHandler {
	return &CategoriaHandler{categoriaService: categoriaService}
}

var _ port.CategoriaHandler = (*CategoriaHandler)(nil)
