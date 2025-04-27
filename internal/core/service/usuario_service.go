package service

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/port"
	"time"
)

type UsuarioService struct {
	usuarioRepository port.UsuarioRepository
}

func (u UsuarioService) ListarUsuarios(ctx context.Context) (*[]domain.UsuarioInfo, error) {
	return u.usuarioRepository.ListarUsuarios(ctx)
}

func (u UsuarioService) ModificarEstadoUsuario(ctx context.Context, usuarioId *int) error {
	return u.usuarioRepository.ModificarEstadoUsuario(ctx, usuarioId)
}

func (u UsuarioService) ModificarUsuario(ctx context.Context, usuarioId *int, usuarioRequest *domain.UsuarioRequest) error {
	if usuarioRequest.DeletedAt != nil {
		*usuarioRequest.DeletedAt = time.Now()
	}
	return u.usuarioRepository.ModificarUsuario(ctx, usuarioId, usuarioRequest)
}

func (u UsuarioService) ObtenerUsuarioDetalle(ctx context.Context, usuarioId *int) (*domain.UsuarioDetalle, error) {
	return u.usuarioRepository.ObtenerUsuarioDetalle(ctx, usuarioId)
}

func (u UsuarioService) RegistrarUsuario(ctx context.Context, usuarioRequest *domain.UsuarioRequest) (*domain.UsuarioDetalle, error) {
	return u.usuarioRepository.RegistrarUsuario(ctx, usuarioRequest)
}

func NewUsuarioService(usuarioRepository port.UsuarioRepository) *UsuarioService {
	return &UsuarioService{usuarioRepository: usuarioRepository}
}

var _ port.UsuarioService = (*UsuarioService)(nil)
