package port

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"mime/multipart"
)

type ProductoRepository interface {
	RegistrarProducto(ctx context.Context, request *domain.ProductRequest, filesHeader *[]*multipart.FileHeader) error
	ModificarProducto(ctx context.Context, id *uuid.UUID, request *domain.ProductRequest, filesHeader *[]*multipart.FileHeader) error
	ListarProductos(ctx context.Context) (*[]domain.ProductoInfo, error)
	ListarUnidadesMedida(ctx context.Context) (*[]domain.UnidadMedida, error)
	ListarFormasFarmaceuticas(ctx context.Context) (*[]domain.FormaFarmaceutica, error)
	HabilitarProducto(ctx context.Context, id *uuid.UUID) error
	DeshabilitarProducto(ctx context.Context, id *uuid.UUID) error
	ObtenerProductoById(ctx context.Context, id *uuid.UUID) (*domain.ProductoDetail, error)
}

type ProductoService interface {
	RegistrarProducto(ctx context.Context, request *domain.ProductRequest, filesHeader *[]*multipart.FileHeader) error
	ModificarProducto(ctx context.Context, id *uuid.UUID, request *domain.ProductRequest, filesHeader *[]*multipart.FileHeader) error
	ListarProductos(ctx context.Context) (*[]domain.ProductoInfo, error)
	ListarUnidadesMedida(ctx context.Context) (*[]domain.UnidadMedida, error)
	ListarFormasFarmaceuticas(ctx context.Context) (*[]domain.FormaFarmaceutica, error)
	HabilitarProducto(ctx context.Context, id *uuid.UUID) error
	DeshabilitarProducto(ctx context.Context, id *uuid.UUID) error
	ObtenerProductoById(ctx context.Context, id *uuid.UUID) (*domain.ProductoDetail, error)
}

type ProductoHandler interface {
	RegistrarProducto(c *fiber.Ctx) error
	ModificarProducto(c *fiber.Ctx) error
	ListarProductos(c *fiber.Ctx) error
	ListarUnidadesMedida(c *fiber.Ctx) error
	ListarFormasFarmaceuticas(c *fiber.Ctx) error
	HabilitarProducto(c *fiber.Ctx) error
	DeshabilitarProducto(c *fiber.Ctx) error
	ObtenerProductoById(c *fiber.Ctx) error
}
