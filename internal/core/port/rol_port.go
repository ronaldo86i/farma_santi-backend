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
	ModificarEstadoRol(ctx context.Context, id *int) error
	ModificarRol(ctx context.Context, id *int, rolRequestUpdate *domain.RolRequestUpdate) error
}

type RolService interface {
	ListarRoles(ctx context.Context) (*[]domain.Rol, error)
	ObtenerRolById(ctx context.Context, id *int) (*domain.Rol, error)
	RegistrarRol(ctx context.Context, rolRequest *domain.RolRequest) error
	ModificarEstadoRol(ctx context.Context, id *int) error
	ModificarRol(ctx context.Context, id *int, rolRequestUpdate *domain.RolRequestUpdate) error
}

type RolHandler interface {
	ListarRoles(c *fiber.Ctx) error
	ObtenerRolById(c *fiber.Ctx) error
	RegistrarRol(c *fiber.Ctx) error
	ModificarEstadoRol(c *fiber.Ctx) error
	ModificarRol(c *fiber.Ctx) error
}
