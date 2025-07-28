package handler

import (
	"errors"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type ProductoHandler struct {
	productoService port.ProductoService
}

func (p ProductoHandler) ObtenerProductoById(c *fiber.Ctx) error {
	// Obtener id del producto del parámetro
	productoIdParam := c.Params("productoId")
	productoId, err := uuid.Parse(productoIdParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(util.NewMessage("Formato de id no válido"))
	}
	producto, err := p.productoService.ObtenerProductoById(c.UserContext(), &productoId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}
	return c.JSON(&producto)
}

func (p ProductoHandler) HabilitarProducto(c *fiber.Ctx) error {
	// Obtener id del producto del parámetro
	productoIdParam := c.Params("productoId")
	productoId, err := uuid.Parse(productoIdParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(util.NewMessage("Formato de id no válido"))
	}
	err = p.productoService.HabilitarProducto(c.UserContext(), &productoId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}
	return c.JSON(util.NewMessage("Producto actualizado correctamente"))
}

func (p ProductoHandler) DeshabilitarProducto(c *fiber.Ctx) error {
	// Obtener id del producto del parámetro
	productoIdParam := c.Params("productoId")
	productoId, err := uuid.Parse(productoIdParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(util.NewMessage("Formato de id no válido"))
	}
	err = p.productoService.DeshabilitarProducto(c.UserContext(), &productoId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}
	return c.JSON(util.NewMessage("Producto actualizado correctamente"))
}

func (p ProductoHandler) ModificarProducto(c *fiber.Ctx) error {
	var productoRequest domain.ProductRequest
	if err := json.Unmarshal([]byte(c.FormValue("body")), &productoRequest); err != nil {
		log.Print(c.FormValue("body"))
		log.Println("Error al deserializar body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(util.NewMessage("Error al leer el formulario"))
	}
	// Obtener id del producto del parámetro
	productoIdParam := c.Params("productoId")
	productoId, err := uuid.Parse(productoIdParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(util.NewMessage("Formato de id no válido"))
	}
	// Obtener el formulario multipart para las imágenes
	form, err := c.MultipartForm()
	if err != nil {
		log.Println("Error al leer el formulario multipart:", err)
		return c.Status(fiber.StatusBadRequest).JSON(util.NewMessage("Error al leer el formulario multipart"))
	}

	files := form.File["images"]
	err = p.productoService.ModificarProducto(c.UserContext(), &productoId, &productoRequest, &files)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}
	return c.JSON(util.NewMessage("Producto modifica correctamente"))
}

func (p ProductoHandler) ListarUnidadesMedida(c *fiber.Ctx) error {
	list, err := p.productoService.ListarUnidadesMedida(c.UserContext())
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

func (p ProductoHandler) ListarFormasFarmaceuticas(c *fiber.Ctx) error {
	list, err := p.productoService.ListarFormasFarmaceuticas(c.UserContext())
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

func (p ProductoHandler) ObtenerListaProductos(c *fiber.Ctx) error {
	list, err := p.productoService.ObtenerListaProductos(c.UserContext())
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

func (p ProductoHandler) RegistrarProducto(c *fiber.Ctx) error {
	var productoRequest domain.ProductRequest
	if err := json.Unmarshal([]byte(c.FormValue("body")), &productoRequest); err != nil {
		log.Print(c.FormValue("body"))
		log.Println("Error al deserializar body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(util.NewMessage("Error al leer el formulario"))
	}
	// Obtener el formulario multipart para las imágenes
	form, err := c.MultipartForm()
	if err != nil {
		log.Println("Error al leer el formulario multipart:", err)
		return c.Status(fiber.StatusBadRequest).JSON(util.NewMessage("Error al leer el formulario multipart"))
	}

	files := form.File["images"]

	err = p.productoService.RegistrarProducto(c.UserContext(), &productoRequest, &files)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerErrorGeneric()
	}
	return c.Status(http.StatusCreated).JSON(util.NewMessage("Producto registrado correctamente"))
}

func NewProductoHandler(productoService port.ProductoService) *ProductoHandler {
	return &ProductoHandler{productoService}
}

var _ port.ProductoHandler = (*ProductoHandler)(nil)
