package service

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
)

type UsuarioService struct {
	usuarioRepository port.UsuarioRepository
}

func (u UsuarioService) RestablecerPassword(ctx context.Context, usuarioId *int, password *domain.UsuarioResetPassword) (*domain.UsuarioDetail, error) {
	return u.usuarioRepository.RestablecerPassword(ctx, usuarioId, password)
}

func (u UsuarioService) HabilitarUsuarioById(ctx context.Context, usuarioId *int) error {
	return u.usuarioRepository.HabilitarUsuarioById(ctx, usuarioId)
}

func (u UsuarioService) DeshabilitarUsuarioById(ctx context.Context, usuarioId *int) error {
	usuario, err := u.usuarioRepository.ObtenerUsuarioDetalle(ctx, usuarioId)
	if err != nil {
		return err
	}
	if usuario.Username == "admin" {
		return datatype.NewBadRequestError("No permitido")
	}
	return u.usuarioRepository.DeshabilitarUsuarioById(ctx, usuarioId)
}

func (u UsuarioService) ObtenerUsuarioDetalleByToken(ctx context.Context, token *string) (*domain.UsuarioDetail, error) {
	claims, err := util.Token.VerifyToken(*token)
	if err != nil {
		return nil, err
	}
	username, ok := claims["username"].(string)
	if !ok {
		return nil, datatype.NewNotFoundError("Usuario no encontrado")
	}
	return u.usuarioRepository.ObtenerUsuarioDetalleByUsername(ctx, &username)
}

func (u UsuarioService) ListarUsuarios(ctx context.Context, filtros map[string]string) (*[]domain.UsuarioInfo, error) {
	return u.usuarioRepository.ListarUsuarios(ctx, filtros)
}

func (u UsuarioService) ModificarUsuario(ctx context.Context, usuarioId *int, usuarioRequest *domain.UsuarioRequest) error {
	return u.usuarioRepository.ModificarUsuario(ctx, usuarioId, usuarioRequest)
}

func (u UsuarioService) ObtenerUsuarioDetalle(ctx context.Context, usuarioId *int) (*domain.UsuarioDetail, error) {
	return u.usuarioRepository.ObtenerUsuarioDetalle(ctx, usuarioId)
}

func (u UsuarioService) RegistrarUsuario(ctx context.Context, usuarioRequest *domain.UsuarioRequest) (*domain.UsuarioDetail, error) {
	return u.usuarioRepository.RegistrarUsuario(ctx, usuarioRequest)
}

func NewUsuarioService(usuarioRepository port.UsuarioRepository) *UsuarioService {
	return &UsuarioService{usuarioRepository: usuarioRepository}
}

var _ port.UsuarioService = (*UsuarioService)(nil)
