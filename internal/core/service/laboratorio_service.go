package service

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/port"
)

type LaboratorioService struct {
	laboratorioRepository port.LaboratorioRepository
}

func (l LaboratorioService) ListarLaboratoriosDisponibles(ctx context.Context) (*[]domain.LaboratorioInfo, error) {
	return l.laboratorioRepository.ListarLaboratoriosDisponibles(ctx)
}

func (l LaboratorioService) ListarLaboratorios(ctx context.Context) (*[]domain.LaboratorioInfo, error) {
	return l.laboratorioRepository.ListarLaboratorios(ctx)
}

func (l LaboratorioService) ObtenerLaboratorioById(ctx context.Context, id *int) (*domain.LaboratorioDetail, error) {
	return l.laboratorioRepository.ObtenerLaboratorioById(ctx, id)
}

func (l LaboratorioService) RegistrarLaboratorio(ctx context.Context, laboratorioRequest *domain.LaboratorioRequest) error {
	return l.laboratorioRepository.RegistrarLaboratorio(ctx, laboratorioRequest)
}

func (l LaboratorioService) ModificarLaboratorio(ctx context.Context, id *int, laboratorioRequest *domain.LaboratorioRequest) error {
	return l.laboratorioRepository.ModificarLaboratorio(ctx, id, laboratorioRequest)
}

func (l LaboratorioService) HabilitarLaboratorio(ctx context.Context, id *int) error {
	return l.laboratorioRepository.HabilitarLaboratorio(ctx, id)
}

func (l LaboratorioService) DeshabilitarLaboratorio(ctx context.Context, id *int) error {
	return l.laboratorioRepository.DeshabilitarLaboratorio(ctx, id)
}

func NewLaboratorioService(laboratorioRepository port.LaboratorioRepository) *LaboratorioService {
	return &LaboratorioService{laboratorioRepository}
}

var _ port.LaboratorioService = (*LaboratorioService)(nil)
