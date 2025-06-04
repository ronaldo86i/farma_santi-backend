package handler

import (
	"errors"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"log"
)

type ProductoHandler struct {
	productoService port.ProductoService
}

func (p ProductoHandler) ListarUnidadesMedida(c *fiber.Ctx) error {
	list, err := p.productoService.ListarUnidadesMedida(c.UserContext())
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

func (p ProductoHandler) ListarFormasFarmaceuticas(c *fiber.Ctx) error {
	list, err := p.productoService.ListarFormasFarmaceuticas(c.UserContext())
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

func (p ProductoHandler) ListarProductos(c *fiber.Ctx) error {
	list, err := p.productoService.ListarProductos(c.UserContext())
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
	// Validar las imágenes
	files := form.File["images"]
	//if len(files) == 0 {
	//	log.Println("No se encontraron archivos con la clave 'images'")
	//	return c.Status(fiber.StatusBadRequest).JSON(util.NewMessage("No se encontraron archivos con la clave 'images'"))
	//}
	err = p.productoService.RegistrarProducto(c.UserContext(), &productoRequest, &files)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return datatype.NewInternalServerError()
	}
	return c.JSON(util.NewMessage("Producto registrado"))
}

func NewProductoHandler(productoService port.ProductoService) *ProductoHandler {
	return &ProductoHandler{productoService}
}

var _ port.ProductoHandler = (*ProductoHandler)(nil)
