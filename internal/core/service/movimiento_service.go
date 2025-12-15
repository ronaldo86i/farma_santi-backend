package service

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/port"
)

type MovimientoService struct {
	movimientoRepository port.MovimientoRepository
}

func (m MovimientoService) ObtenerMovimientosKardex(ctx context.Context, filtros map[string]string) (*[]domain.MovimientoKardex, error) {
	return m.movimientoRepository.ObtenerMovimientosKardex(ctx, filtros)
}

func (m MovimientoService) ObtenerListaMovimientos(ctx context.Context, filtros map[string]string) (*[]domain.MovimientoInfo, error) {
	return m.movimientoRepository.ObtenerListaMovimientos(ctx, filtros)
}

func NewMovimientoService(movimientoRepository port.MovimientoRepository) *MovimientoService {
	return &MovimientoService{movimientoRepository: movimientoRepository}
}

var _ port.MovimientoService = (*MovimientoService)(nil)
