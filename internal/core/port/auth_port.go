package port

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	ObtenerTokenByCredencial(ctx context.Context, credentials *domain.LoginRequest) (*domain.TokenResponse, error)
}

type AuthHandler interface {
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	RefreshOrVerify(c *fiber.Ctx) error
}
