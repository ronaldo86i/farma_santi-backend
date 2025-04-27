package service

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/port"
)

type ProveedorService struct {
	proveedorRepository port.ProveedorRepository
}

func (p ProveedorService) RegistrarProveedor(ctx context.Context, request *domain.ProveedorRequest) error {
	return p.proveedorRepository.RegistrarProveedor(ctx, request)
}

func (p ProveedorService) ObtenerProveedorById(ctx context.Context, id *int) (*domain.ProveedorDetail, error) {
	return p.proveedorRepository.ObtenerProveedorById(ctx, id)
}

func (p ProveedorService) ListarProveedores(ctx context.Context) (*[]domain.ProveedorInfo, error) {
	return p.proveedorRepository.ListarProveedores(ctx)
}

func (p ProveedorService) ModificarProveedor(ctx context.Context, id *int, request *domain.ProveedorRequest) error {
	return p.proveedorRepository.ModificarProveedor(ctx, id, request)
}

func (p ProveedorService) ModificarEstadoProveedor(ctx context.Context, id *int) error {
	return p.proveedorRepository.ModificarEstadoProveedor(ctx, id)
}

func NewProveedorService(proveedorRepository port.ProveedorRepository) *ProveedorService {
	return &ProveedorService{proveedorRepository}
}

var _ port.ProveedorService = (*ProveedorService)(nil)
