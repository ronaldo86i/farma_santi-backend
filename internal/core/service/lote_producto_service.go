package service

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/port"
	"github.com/google/uuid"
)

type LoteProductoService struct {
	loteProductoRepository port.LoteProductoRepository
}

func (l LoteProductoService) ActualizarLotesVencidos(ctx context.Context) error {
	return l.loteProductoRepository.ActualizarLotesVencidos(ctx)
}

func (l LoteProductoService) ListarLotesProductosByProductoId(ctx context.Context, productoId *uuid.UUID) (*[]domain.LoteProductoSimple, error) {
	return l.loteProductoRepository.ListarLotesProductosByProductoId(ctx, productoId)
}

func (l LoteProductoService) ModificarLoteProducto(ctx context.Context, id *int, request *domain.LoteProductoRequest) error {
	return l.loteProductoRepository.ModificarLoteProducto(ctx, id, request)
}

func (l LoteProductoService) ListarLotesProductos(ctx context.Context) (*[]domain.LoteProductoInfo, error) {
	return l.loteProductoRepository.ListarLotesProductos(ctx)
}

func (l LoteProductoService) RegistrarLoteProducto(ctx context.Context, request *domain.LoteProductoRequest) error {
	return l.loteProductoRepository.RegistrarLoteProducto(ctx, request)
}

func (l LoteProductoService) ObtenerLoteProductoById(ctx context.Context, id *int) (*domain.LoteProductoDetail, error) {
	return l.loteProductoRepository.ObtenerLoteProductoById(ctx, id)
}

func NewLoteProductoService(loteProductoRepository port.LoteProductoRepository) *LoteProductoService {
	return &LoteProductoService{loteProductoRepository: loteProductoRepository}
}

var _ port.LoteProductoService = (*LoteProductoService)(nil)
