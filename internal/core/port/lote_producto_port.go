package port

import (
	"context"
	"farma-santi_backend/internal/core/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LoteProductoRepository interface {
	ObtenerListaLotesProductos(ctx context.Context, filtros map[string]string) (*[]domain.LoteProductoInfo, error)
	RegistrarLoteProducto(ctx context.Context, request *domain.LoteProductoRequest) error
	ModificarLoteProducto(ctx context.Context, id *int, request *domain.LoteProductoRequest) error
	ObtenerLoteProductoById(ctx context.Context, id *int) (*domain.LoteProductoDetail, error)
	ListarLotesProductosByProductoId(ctx context.Context, productoId *uuid.UUID) (*[]domain.LoteProductoSimple, error)
	ActualizarLotesVencidos(ctx context.Context) error
}

type LoteProductoService interface {
	ObtenerListaLotesProductos(ctx context.Context, filtros map[string]string) (*[]domain.LoteProductoInfo, error)
	RegistrarLoteProducto(ctx context.Context, request *domain.LoteProductoRequest) error
	ModificarLoteProducto(ctx context.Context, id *int, request *domain.LoteProductoRequest) error
	ObtenerLoteProductoById(ctx context.Context, id *int) (*domain.LoteProductoDetail, error)
	ListarLotesProductosByProductoId(ctx context.Context, productoId *uuid.UUID) (*[]domain.LoteProductoSimple, error)
	ActualizarLotesVencidos(ctx context.Context) error
}

type LoteProductoHandler interface {
	ObtenerListaLotesProductos(c *fiber.Ctx) error
	RegistrarLoteProducto(c *fiber.Ctx) error
	ModificarLoteProducto(c *fiber.Ctx) error
	ObtenerLoteProductoById(c *fiber.Ctx) error
	ListarLotesProductosByProductoId(c *fiber.Ctx) error
}
