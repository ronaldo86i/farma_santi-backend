package service

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/port"
)

type MovimientoService struct {
	movimientoRepository port.MovimientoRepository
}

func (m MovimientoService) ObtenerListaMovimientos(ctx context.Context) (*[]domain.MovimientoInfo, error) {
	return m.movimientoRepository.ObtenerListaMovimientos(ctx)
}

func NewMovimientoService(movimientoRepository port.MovimientoRepository) *MovimientoService {
	return &MovimientoService{movimientoRepository: movimientoRepository}
}

var _ port.MovimientoService = (*MovimientoService)(nil)
