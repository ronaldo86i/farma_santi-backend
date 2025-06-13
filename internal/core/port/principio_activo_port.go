package port

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
)

type PrincipioActivoRepository interface {
	RegistrarPrincipioActivo(ctx context.Context, request *domain.PrincipioActivoRequest) error
	ModificarPrincipioActivo(ctx context.Context, id *int, request *domain.PrincipioActivoRequest) error
	ListarPrincipioActivo(ctx context.Context) (*[]domain.PrincipioActivoInfo, error)
	ObtenerPrincipioActivoById(ctx context.Context, id *int) (*domain.PrincipioActivoDetail, error)
}

type PrincipioActivoService interface {
	RegistrarPrincipioActivo(ctx context.Context, request *domain.PrincipioActivoRequest) error
	ModificarPrincipioActivo(ctx context.Context, id *int, request *domain.PrincipioActivoRequest) error
	ListarPrincipioActivo(ctx context.Context) (*[]domain.PrincipioActivoInfo, error)
	ObtenerPrincipioActivoById(ctx context.Context, id *int) (*domain.PrincipioActivoDetail, error)
}

type PrincipioActivoHandler interface {
	RegistrarPrincipioActivo(c *fiber.Ctx) error
	ModificarPrincipioActivo(c *fiber.Ctx) error
	ListarPrincipioActivo(c *fiber.Ctx) error
	ObtenerPrincipioActivoById(c *fiber.Ctx) error
}
