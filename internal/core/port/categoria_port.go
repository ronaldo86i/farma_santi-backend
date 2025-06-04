package port

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
)

type CategoriaRepository interface {
	ListarCategoriasDisponibles(ctx context.Context) (*[]domain.Categoria, error)
	RegistrarCategoria(ctx context.Context, categoriaRequest *domain.CategoriaRequest) error
	ObtenerCategoriaById(ctx context.Context, categoriaId *int) (*domain.Categoria, error)
	ListarCategorias(ctx context.Context) (*[]domain.Categoria, error)
	ModificarCategoria(ctx context.Context, categoriaId *int, categoriaRequest *domain.CategoriaRequest) error
	HabilitarCategoria(ctx context.Context, categoriaId *int) error
	DeshabilitarCategoria(ctx context.Context, categoriaId *int) error
}

type CategoriaService interface {
	ListarCategoriasDisponibles(ctx context.Context) (*[]domain.Categoria, error)
	RegistrarCategoria(ctx context.Context, categoriaRequest *domain.CategoriaRequest) error
	ObtenerCategoriaById(ctx context.Context, categoriaId *int) (*domain.Categoria, error)
	ListarCategorias(ctx context.Context) (*[]domain.Categoria, error)
	ModificarCategoria(ctx context.Context, categoriaId *int, categoriaRequest *domain.CategoriaRequest) error
	HabilitarCategoria(ctx context.Context, categoriaId *int) error
	DeshabilitarCategoria(ctx context.Context, categoriaId *int) error
}

type CategoriaHandler interface {
	ListarCategoriasDisponibles(c *fiber.Ctx) error
	RegistrarCategoria(c *fiber.Ctx) error
	ObtenerCategoriaById(c *fiber.Ctx) error
	ListarCategorias(c *fiber.Ctx) error
	ModificarCategoria(c *fiber.Ctx) error
	HabilitarCategoria(c *fiber.Ctx) error
	DeshabilitarCategoria(c *fiber.Ctx) error
}
