package service

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/port"
)

type PresentacionService struct {
	presentacionRepository port.PresentacionRepository
}

func (p PresentacionService) ObtenerListaPresentaciones(ctx context.Context) (*[]domain.Presentacion, error) {
	return p.presentacionRepository.ObtenerListaPresentaciones(ctx)
}

func NewPresentacionService(presentacionRepository port.PresentacionRepository) *PresentacionService {
	return &PresentacionService{presentacionRepository: presentacionRepository}
}

var _ port.PresentacionService = (*PresentacionService)(nil)
