package port

import (
	"context"
	"farma-santi_backend/internal/core/domain"

	"github.com/gofiber/fiber/v2"
)

type ClienteRepository interface {
	ObtenerListaClientes(ctx context.Context, filtros map[string]string) (*[]domain.ClienteInfo, error)
	ObtenerClienteById(ctx context.Context, id *int) (*domain.ClienteDetail, error)
	RegistrarCliente(ctx context.Context, request *domain.ClienteRequest) (*int, error)
	ModificarClienteById(ctx context.Context, id *int, request *domain.ClienteRequest) error
	HabilitarCliente(ctx context.Context, id *int) error
	DeshabilitarCliente(ctx context.Context, id *int) error
}

type ClienteService interface {
	ObtenerListaClientes(ctx context.Context, filtros map[string]string) (*[]domain.ClienteInfo, error)
	ObtenerClienteById(ctx context.Context, id *int) (*domain.ClienteDetail, error)
	RegistrarCliente(ctx context.Context, request *domain.ClienteRequest) (*int, error)
	ModificarClienteById(ctx context.Context, id *int, request *domain.ClienteRequest) error
	HabilitarCliente(ctx context.Context, id *int) error
	DeshabilitarCliente(ctx context.Context, id *int) error
}

type ClienteHandler interface {
	ObtenerListaClientes(c *fiber.Ctx) error
	ObtenerClienteById(c *fiber.Ctx) error
	RegistrarCliente(c *fiber.Ctx) error
	ModificarClienteById(c *fiber.Ctx) error
	HabilitarCliente(c *fiber.Ctx) error
	DeshabilitarCliente(c *fiber.Ctx) error
}
