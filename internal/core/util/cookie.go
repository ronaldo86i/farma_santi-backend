package util

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

func SetCookie(c *fiber.Ctx, name string, value string, duration time.Duration, httpOnly, secure bool, timeNow time.Time) {
	if name == "" {
		return
	}

	exp := timeNow.Add(duration)

	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    value,
		Expires:  exp,
		MaxAge:   int(duration.Seconds()),
		HTTPOnly: httpOnly,
		Secure:   secure,
		SameSite: "Lax",
	})
}

// DeleteCookie elimina una cookie con nombre especificado
func DeleteCookie(c *fiber.Ctx, name string, httpOnly bool) {
	expired := time.Now().Add(-2 * time.Hour)

	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    "",
		Expires:  expired,
		MaxAge:   -1,
		Path:     "/",
		HTTPOnly: httpOnly,
		Secure:   false,
		SameSite: "Lax",
	})
}
