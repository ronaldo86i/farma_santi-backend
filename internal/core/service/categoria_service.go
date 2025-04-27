package service

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/port"
	"time"
)

type CategoriaService struct {
	categoriaRepository port.CategoriaRepository
}

func (c CategoriaService) ObtenerCategoriaById(ctx context.Context, categoriaId *int) (*domain.Categoria, error) {
	return c.categoriaRepository.ObtenerCategoriaById(ctx, categoriaId)
}

func (c CategoriaService) ListarCategorias(ctx context.Context) (*[]domain.Categoria, error) {
	return c.categoriaRepository.ListarCategorias(ctx)
}

func (c CategoriaService) ModificarEstadoCategoria(ctx context.Context, categoriaId *int) error {
	return c.categoriaRepository.ModificarEstadoCategoria(ctx, categoriaId)
}

func (c CategoriaService) ModificarCategoria(ctx context.Context, categoriaId *int, categoriaRequest *domain.CategoriaRequest) error {
	if categoriaRequest.DeletedAt != nil {
		*categoriaRequest.DeletedAt = time.Now()
	}
	return c.categoriaRepository.ModificarCategoria(ctx, categoriaId, categoriaRequest)
}

func (c CategoriaService) RegistrarCategoria(ctx context.Context, categoriaRequest *domain.CategoriaRequest) error {
	return c.categoriaRepository.RegistrarCategoria(ctx, categoriaRequest)
}

func NewCategoriaService(categoriaRepository port.CategoriaRepository) *CategoriaService {
	return &CategoriaService{categoriaRepository: categoriaRepository}
}

var _ port.CategoriaService = (*CategoriaService)(nil)
