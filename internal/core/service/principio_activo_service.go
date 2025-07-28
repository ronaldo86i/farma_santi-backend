package service

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/port"
	"strings"
)

type PrincipioActivoService struct {
	principioActivoRepository port.PrincipioActivoRepository
}

func (p PrincipioActivoService) RegistrarPrincipioActivo(ctx context.Context, request *domain.PrincipioActivoRequest) (*int, error) {
	request.Nombre = strings.TrimSpace(request.Nombre)
	request.Nombre = strings.ToUpper(request.Nombre)
	return p.principioActivoRepository.RegistrarPrincipioActivo(ctx, request)
}

func (p PrincipioActivoService) ModificarPrincipioActivo(ctx context.Context, id *int, request *domain.PrincipioActivoRequest) error {
	request.Nombre = strings.TrimSpace(request.Nombre)
	request.Nombre = strings.ToUpper(request.Nombre)
	return p.principioActivoRepository.ModificarPrincipioActivo(ctx, id, request)
}

func (p PrincipioActivoService) ListarPrincipioActivo(ctx context.Context) (*[]domain.PrincipioActivoInfo, error) {
	return p.principioActivoRepository.ListarPrincipioActivo(ctx)
}

func (p PrincipioActivoService) ObtenerPrincipioActivoById(ctx context.Context, id *int) (*domain.PrincipioActivoDetail, error) {
	return p.principioActivoRepository.ObtenerPrincipioActivoById(ctx, id)
}

func NewPrincipioActivoService(activoRepository port.PrincipioActivoRepository) *PrincipioActivoService {
	return &PrincipioActivoService{principioActivoRepository: activoRepository}
}

var _ port.PrincipioActivoService = (*PrincipioActivoService)(nil)
