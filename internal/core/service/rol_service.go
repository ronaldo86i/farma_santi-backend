package service

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/port"
	"strings"
)

type RolService struct {
	rolRepository port.RolRepository
}

func (r RolService) HabilitarRol(ctx context.Context, id *int) error {
	return r.rolRepository.HabilitarRol(ctx, id)
}

func (r RolService) DeshabilitarRol(ctx context.Context, id *int) error {
	return r.rolRepository.DeshabilitarRol(ctx, id)
}

func (r RolService) ModificarRol(ctx context.Context, id *int, rolRequestUpdate *domain.RolRequest) error {
	rolRequestUpdate.Nombre = strings.ToUpper(rolRequestUpdate.Nombre)
	return r.rolRepository.ModificarRol(ctx, id, rolRequestUpdate)
}

func (r RolService) RegistrarRol(ctx context.Context, rolRequest *domain.RolRequest) error {
	rolRequest.Nombre = strings.ToUpper(rolRequest.Nombre)
	return r.rolRepository.RegistrarRol(ctx, rolRequest)
}

func (r RolService) ListarRoles(ctx context.Context) (*[]domain.Rol, error) {
	return r.rolRepository.ListarRoles(ctx)
}

func (r RolService) ObtenerRolById(ctx context.Context, id *int) (*domain.Rol, error) {
	return r.rolRepository.ObtenerRolById(ctx, id)
}

func NewRolService(rolRepository port.RolRepository) *RolService {
	return &RolService{rolRepository}
}

var _ port.RolService = (*RolService)(nil)
