package port

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
)

type CategoriaRepository interface {
	RegistrarCategoria(ctx context.Context, categoriaRequest *domain.CategoriaRequest) error
	ObtenerCategoriaById(ctx context.Context, categoriaId *int) (*domain.Categoria, error)
	ListarCategorias(ctx context.Context) (*[]domain.Categoria, error)
	ModificarEstadoCategoria(ctx context.Context, categoriaId *int) error
	ModificarCategoria(ctx context.Context, categoriaId *int, categoriaRequest *domain.CategoriaRequest) error
}

type CategoriaService interface {
	RegistrarCategoria(ctx context.Context, categoriaRequest *domain.CategoriaRequest) error
	ObtenerCategoriaById(ctx context.Context, categoriaId *int) (*domain.Categoria, error)
	ListarCategorias(ctx context.Context) (*[]domain.Categoria, error)
	ModificarEstadoCategoria(ctx context.Context, categoriaId *int) error
	ModificarCategoria(ctx context.Context, categoriaId *int, categoriaRequest *domain.CategoriaRequest) error
}

type CategoriaHandler interface {
	RegistrarCategoria(c *fiber.Ctx) error
	ObtenerCategoriaById(c *fiber.Ctx) error
	ListarCategorias(c *fiber.Ctx) error
	ModificarEstadoCategoria(c *fiber.Ctx) error
	ModificarCategoria(c *fiber.Ctx) error
}
