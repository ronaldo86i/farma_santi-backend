package port

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
)

type LaboratorioRepository interface {
	ListarLaboratoriosDisponibles(ctx context.Context) (*[]domain.LaboratorioInfo, error)
	ListarLaboratorios(ctx context.Context) (*[]domain.LaboratorioInfo, error)
	ObtenerLaboratorioById(ctx context.Context, id *int) (*domain.LaboratorioDetail, error)
	RegistrarLaboratorio(ctx context.Context, laboratorioRequest *domain.LaboratorioRequest) error
	ModificarLaboratorio(ctx context.Context, id *int, laboratorioRequest *domain.LaboratorioRequest) error
	HabilitarLaboratorio(ctx context.Context, id *int) error
	DeshabilitarLaboratorio(ctx context.Context, id *int) error
}

type LaboratorioService interface {
	ListarLaboratoriosDisponibles(ctx context.Context) (*[]domain.LaboratorioInfo, error)
	ListarLaboratorios(ctx context.Context) (*[]domain.LaboratorioInfo, error)
	ObtenerLaboratorioById(ctx context.Context, id *int) (*domain.LaboratorioDetail, error)
	RegistrarLaboratorio(ctx context.Context, laboratorioRequest *domain.LaboratorioRequest) error
	ModificarLaboratorio(ctx context.Context, id *int, laboratorioRequest *domain.LaboratorioRequest) error
	HabilitarLaboratorio(ctx context.Context, id *int) error
	DeshabilitarLaboratorio(ctx context.Context, id *int) error
}

type LaboratorioHandler interface {
	ListarLaboratoriosDisponibles(c *fiber.Ctx) error
	ListarLaboratorios(c *fiber.Ctx) error
	ObtenerLaboratorioById(c *fiber.Ctx) error
	RegistrarLaboratorio(c *fiber.Ctx) error
	ModificarLaboratorio(c *fiber.Ctx) error
	HabilitarLaboratorio(c *fiber.Ctx) error
	DeshabilitarLaboratorio(c *fiber.Ctx) error
}
