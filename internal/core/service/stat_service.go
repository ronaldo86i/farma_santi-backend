package service

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/port"
)

type StatService struct {
	statRepository port.StatRepository
}

func (s StatService) ObtenerEstadisticasDashboard(ctx context.Context) (*domain.DashboardStats, error) {
	return s.statRepository.ObtenerEstadisticasDashboard(ctx)
}

func (s StatService) ObtenerTopProductosVendidos(ctx context.Context, filtros map[string]string) (*[]domain.ProductoStat, error) {
	return s.statRepository.ObtenerTopProductosVendidos(ctx, filtros)
}

func NewStatService(statRepository port.StatRepository) *StatService {
	return &StatService{statRepository: statRepository}
}

var _ port.StatService = (*StatService)(nil)
