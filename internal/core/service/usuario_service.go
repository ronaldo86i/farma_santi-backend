package service

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
	"net/http"
	"time"
)

type UsuarioService struct {
	usuarioRepository port.UsuarioRepository
}

func (u UsuarioService) RestablecerPassword(ctx context.Context, usuarioId *int) (*domain.UsuarioDetail, error) {
	return u.usuarioRepository.RestablecerPassword(ctx, usuarioId)
}

func (u UsuarioService) HabilitarUsuarioById(ctx context.Context, usuarioId *int) error {
	return u.usuarioRepository.HabilitarUsuarioById(ctx, usuarioId)
}

func (u UsuarioService) DeshabilitarUsuarioById(ctx context.Context, usuarioId *int) error {
	return u.usuarioRepository.DeshabilitarUsuarioById(ctx, usuarioId)
}

func (u UsuarioService) ObtenerUsuarioDetalleByToken(ctx context.Context, token *string) (*domain.UsuarioDetail, error) {
	claims, err := util.Token.VerifyToken(*token)
	if err != nil {
		return nil, err
	}
	username, ok := claims["username"].(string)
	if !ok {
		return nil, &datatype.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Usuario no encontrado",
		}
	}
	return u.usuarioRepository.ObtenerUsuarioDetalleByUsername(ctx, &username)
}

func (u UsuarioService) ListarUsuarios(ctx context.Context) (*[]domain.UsuarioInfo, error) {
	return u.usuarioRepository.ListarUsuarios(ctx)
}

func (u UsuarioService) ModificarUsuario(ctx context.Context, usuarioId *int, usuarioRequest *domain.UsuarioRequest) error {
	if usuarioRequest.DeletedAt != nil {
		*usuarioRequest.DeletedAt = time.Now()
	}
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
