package port

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
)

type ProductoRepository interface {
	RegistrarProducto(ctx context.Context, request *domain.ProductRequest, filesHeader *[]*multipart.FileHeader) error
	ListarProductos(ctx context.Context) (*[]domain.ProductInfo, error)
	ListarUnidadesMedida(ctx context.Context) (*[]domain.UnidadMedida, error)
	ListarFormasFarmaceuticas(ctx context.Context) (*[]domain.FormaFarmaceutica, error)
}

type ProductoService interface {
	RegistrarProducto(ctx context.Context, request *domain.ProductRequest, filesHeader *[]*multipart.FileHeader) error
	ListarProductos(ctx context.Context) (*[]domain.ProductInfo, error)
	ListarUnidadesMedida(ctx context.Context) (*[]domain.UnidadMedida, error)
	ListarFormasFarmaceuticas(ctx context.Context) (*[]domain.FormaFarmaceutica, error)
}

type ProductoHandler interface {
	RegistrarProducto(c *fiber.Ctx) error
	ListarProductos(c *fiber.Ctx) error
	ListarUnidadesMedida(c *fiber.Ctx) error
	ListarFormasFarmaceuticas(c *fiber.Ctx) error
}
