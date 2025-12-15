package util

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetCookie(c *fiber.Ctx, name string, value string, duration time.Duration, httpOnly, secure bool, timeNow time.Time) {

	domain := os.Getenv("COOKIE_DOMAIN")
	if domain == "" {
		domain = "localhost"
	}
	exp := timeNow.Add(duration)

	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    value,
		Expires:  exp,
		MaxAge:   int(duration.Seconds()),
		HTTPOnly: httpOnly,
		Domain:   domain,
		Secure:   secure,
		SameSite: "Lax",
	})
}

// DeleteCookie elimina una cookie con nombre especificado
func DeleteCookie(c *fiber.Ctx, name string, httpOnly bool) {
	domain := os.Getenv("COOKIE_DOMAIN")
	if domain == "" {
		domain = "localhost"
	}
	expired := time.Now().Add(-2 * time.Hour)

	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    "",
		Expires:  expired,
		MaxAge:   -1,
		Path:     "/",
		HTTPOnly: httpOnly,
		Domain:   domain,
		Secure:   false,
		SameSite: "Lax",
	})
}
