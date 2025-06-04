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

type UsuarioHandler struct {
	usuarioService port.UsuarioService
}

func (u UsuarioHandler) RestablecerPassword(c *fiber.Ctx) error {
	ctx := c.UserContext() // Usa el contexto de la petición
	usuarioId, err := c.ParamsInt("usuarioId")
	if err != nil || usuarioId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del usuario debe ser un número válido mayor a 0"))
	}
	usuario, err := u.usuarioService.RestablecerPassword(ctx, &usuarioId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(util.NewMessageData(&usuario, "Contraseña actualizada correctamente"))
}

func (u UsuarioHandler) HabilitarUsuarioById(c *fiber.Ctx) error {
	ctx := c.UserContext() // Usa el contexto de la petición
	usuarioId, err := c.ParamsInt("usuarioId")
	if err != nil || usuarioId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del usuario debe ser un número válido mayor a 0"))
	}

	err = u.usuarioService.HabilitarUsuarioById(ctx, &usuarioId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(util.NewMessage("Usuario actualizado correctamente"))
}

func (u UsuarioHandler) DeshabilitarUsuarioById(c *fiber.Ctx) error {
	ctx := c.UserContext() // Usa el contexto de la petición
	usuarioId, err := c.ParamsInt("usuarioId")
	if err != nil || usuarioId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del usuario debe ser un número válido mayor a 0"))
	}

	err = u.usuarioService.DeshabilitarUsuarioById(ctx, &usuarioId)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(util.NewMessage("Usuario actualizado correctamente"))
}

func (u UsuarioHandler) ObtenerUsuarioActual(c *fiber.Ctx) error {
	token := c.Cookies("access-token")
	usuarioDetalle, err := u.usuarioService.ObtenerUsuarioDetalleByToken(c.Context(), &token)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(&usuarioDetalle)
}

func (u UsuarioHandler) ListarUsuarios(c *fiber.Ctx) error {

	listaUsuario, err := u.usuarioService.ListarUsuarios(c.Context())
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}
	return c.JSON(&listaUsuario)
}

func (u UsuarioHandler) ModificarUsuario(c *fiber.Ctx) error {
	ctx := c.UserContext() // Usar el contexto

	var usuarioRequest domain.UsuarioRequest
	if err := c.BodyParser(&usuarioRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("Petición inválida: datos incompletos o incorrectos"))
	}

	usuarioId, err := c.ParamsInt("usuarioId")
	if err != nil || usuarioId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del rol debe ser un número válido mayor a 0"))
	}
	err = u.usuarioService.ModificarUsuario(ctx, &usuarioId, &usuarioRequest)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}

	return c.Status(http.StatusAccepted).JSON(util.NewMessage("Usuario actualizado correctamente"))
}

func (u UsuarioHandler) ObtenerUsuarioDetalle(c *fiber.Ctx) error {
	ctx := c.UserContext()
	usuarioId, err := c.ParamsInt("usuarioId", 0)
	if err != nil || usuarioId <= 0 {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("El 'id' del usuario debe ser un número válido mayor a 0"))
	}
	usuarioDetalle, err := u.usuarioService.ObtenerUsuarioDetalle(ctx, &usuarioId)

	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(&usuarioDetalle)
}

func (u UsuarioHandler) RegistrarUsuario(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var usuarioRequest domain.UsuarioRequest
	if err := c.BodyParser(&usuarioRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("Petición inválida: datos incompletos o incorrectos"))
	}

	usuarioDetalle, err := u.usuarioService.RegistrarUsuario(ctx, &usuarioRequest)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}

	return c.Status(http.StatusCreated).JSON(util.NewMessageData(usuarioDetalle, "Usuario creado correctamente"))
}

func NewUsuarioHandler(usuarioService port.UsuarioService) UsuarioHandler {
	return UsuarioHandler{usuarioService}
}

var _ port.UsuarioHandler = (*UsuarioHandler)(nil)
