package middleware

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

const hostnameKey = "fullHostname"

// HostnameMiddleware guarda y registra el hostname completo de la petición
func HostnameMiddleware(c *fiber.Ctx) error {
	fullHostname := fmt.Sprintf("%s://%s", c.Protocol(), c.Hostname())
	// Aquí podrías guardar el hostname en una base de datos, archivo, etc.
	log.Printf("Petición recibida desde host: %s", fullHostname)

	// Guardar también en context.Context estándar (útil si integras con librerías externas)
	ctx := context.WithValue(c.UserContext(), hostnameKey, fullHostname)
	c.SetUserContext(ctx)
	// También puedes guardar el hostname en el contexto para acceder después
	c.Locals("hostname", fullHostname)

	return c.Next()
}
