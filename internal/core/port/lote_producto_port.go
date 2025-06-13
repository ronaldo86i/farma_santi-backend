package port

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
)

type LoteProductoRepository interface {
	ListarLotesProductos(ctx context.Context) (*[]domain.LoteProductoInfo, error)
	RegistrarLoteProducto(ctx context.Context, request *domain.LoteProductoRequest) error
	ModificarLoteProducto(ctx context.Context, id *int, request *domain.LoteProductoRequest) error
	ObtenerLoteProductoById(ctx context.Context, id *int) (*domain.LoteProductoDetail, error)
}

type LoteProductoService interface {
	ListarLotesProductos(ctx context.Context) (*[]domain.LoteProductoInfo, error)
	RegistrarLoteProducto(ctx context.Context, request *domain.LoteProductoRequest) error
	ModificarLoteProducto(ctx context.Context, id *int, request *domain.LoteProductoRequest) error
	ObtenerLoteProductoById(ctx context.Context, id *int) (*domain.LoteProductoDetail, error)
}

type LoteProductoHandler interface {
	ListarLotesProductos(c *fiber.Ctx) error
	RegistrarLoteProducto(c *fiber.Ctx) error
	ModificarLoteProducto(c *fiber.Ctx) error
	ObtenerLoteProductoById(c *fiber.Ctx) error
}
