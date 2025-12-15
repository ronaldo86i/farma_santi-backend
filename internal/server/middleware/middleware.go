package middleware

import (
	"context"
	"errors"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/util"
	"farma-santi_backend/internal/server/setup"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// HostnameMiddleware guarda y registra el hostname completo de la petición
func HostnameMiddleware(c *fiber.Ctx) error {
	fullHostname := fmt.Sprintf("%s://%s", c.Protocol(), c.Hostname())
	log.Printf("Petición recibida desde host: %s", fullHostname)
	// Guardar fullHostname en context
	ctx := context.WithValue(c.UserContext(), util.ContextFullHostnameKey, fullHostname)
	c.SetUserContext(ctx)
	return c.Next()
}

func VerifyUserAdminMiddleware(c *fiber.Ctx) error {
	claimsAccessToken, err := util.Token.VerifyToken(c.Cookies("access-token"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(util.NewMessage("Usuario no autorizado"))
	}

	// Extraer username
	username, ok := claimsAccessToken["username"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(util.NewMessage("Usuario no autorizado"))
	}

	// Extraer userId
	userIdFloat, ok := claimsAccessToken["userId"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(util.NewMessage("Usuario no autorizado"))
	}
	userId := int(userIdFloat)
	// Guardar en el contexto
	ctx := context.WithValue(c.UserContext(), util.ContextUsernameKey, username)
	ctx = context.WithValue(ctx, util.ContextUserIdKey, userId)
	c.SetUserContext(ctx)

	return c.Next()
}

func VerifyRolesMiddleware(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.UserContext().Value(util.ContextUsernameKey).(string)
		user, err := setup.GetDependencies().Repository.Usuario.ObtenerUsuarioDetalleByUsername(c.UserContext(), &username)
		if err != nil {
			log.Print(err.Error())
			var errorResponse *datatype.ErrorResponse
			if errors.As(err, &errorResponse) {
				return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
			}
			return datatype.NewInternalServerErrorGeneric()
		}
		if user.Estado == "Inactivo" {
			return datatype.NewBadRequestError("Usuario no autorizado")
		}
		// Verificar si tiene el rol
		for _, rol := range roles {
			for _, userRole := range user.Roles {
				if rol == userRole.Nombre {
					return c.Next()
				}
			}
		}

		return c.Status(http.StatusUnauthorized).JSON(util.NewMessage("Usuario no autorizado"))
	}
}

func VerifyUsuarioShared(c *fiber.Ctx) error {
	tokenString, err := util.Token.GetToken(c.Get("Authorization"))
	if err != nil {
		log.Println("Error al obtener token", err)
		return c.Status(fiber.StatusUnauthorized).JSON(util.NewMessage("Usuario no autorizado"))
	}

	claimsAccessToken, err := util.Token.VerifyTokenType(tokenString, "access-token-public")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(util.NewMessage("Usuario no autorizado"))
	}

	// Extraer userId
	userId, ok := claimsAccessToken["userId"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(util.NewMessage("Usuario no autorizado"))
	}

	// Guardar en el contexto
	ctx := context.WithValue(c.UserContext(), util.ContextUserIdKey, userId)
	c.SetUserContext(ctx)

	// Guardar en local
	c.Locals(util.ContextUserIdKey, userId)

	return c.Next()
}
