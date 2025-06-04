package port

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
)

type RolRepository interface {
	ListarRoles(ctx context.Context) (*[]domain.Rol, error)
	ObtenerRolById(ctx context.Context, id *int) (*domain.Rol, error)
	RegistrarRol(ctx context.Context, rolRequest *domain.RolRequest) error
	ModificarRol(ctx context.Context, id *int, rolRequest *domain.RolRequest) error
	HabilitarRol(ctx context.Context, id *int) error
	DeshabilitarRol(ctx context.Context, id *int) error
}

type RolService interface {
	ListarRoles(ctx context.Context) (*[]domain.Rol, error)
	ObtenerRolById(ctx context.Context, id *int) (*domain.Rol, error)
	RegistrarRol(ctx context.Context, rolRequest *domain.RolRequest) error
	ModificarRol(ctx context.Context, id *int, rolRequest *domain.RolRequest) error
	HabilitarRol(ctx context.Context, id *int) error
	DeshabilitarRol(ctx context.Context, id *int) error
}

type RolHandler interface {
	ListarRoles(c *fiber.Ctx) error
	ObtenerRolById(c *fiber.Ctx) error
	RegistrarRol(c *fiber.Ctx) error
	ModificarRol(c *fiber.Ctx) error
	HabilitarRol(c *fiber.Ctx) error
	DeshabilitarRol(c *fiber.Ctx) error
}
