package port

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
)

type UsuarioRepository interface {
	HabilitarUsuarioById(ctx context.Context, usuarioId *int) error
	DeshabilitarUsuarioById(ctx context.Context, usuarioId *int) error
	ModificarUsuario(ctx context.Context, usuarioId *int, usuarioRequest *domain.UsuarioRequest) error
	ObtenerUsuario(ctx context.Context, username *string) (*domain.Usuario, error)
	ObtenerUsuarioDetalle(ctx context.Context, usuarioId *int) (*domain.UsuarioDetail, error)
	ObtenerUsuarioDetalleByUsername(ctx context.Context, username *string) (*domain.UsuarioDetail, error)
	RegistrarUsuario(ctx context.Context, usuarioRequest *domain.UsuarioRequest) (*domain.UsuarioDetail, error)
	ListarUsuarios(ctx context.Context) (*[]domain.UsuarioInfo, error)
	RestablecerPassword(ctx context.Context, usuarioId *int) (*domain.UsuarioDetail, error)
}

type UsuarioService interface {
	HabilitarUsuarioById(ctx context.Context, usuarioId *int) error
	DeshabilitarUsuarioById(ctx context.Context, usuarioId *int) error
	ModificarUsuario(ctx context.Context, usuarioId *int, usuarioRequest *domain.UsuarioRequest) error
	ObtenerUsuarioDetalle(ctx context.Context, usuarioId *int) (*domain.UsuarioDetail, error)
	ObtenerUsuarioDetalleByToken(ctx context.Context, token *string) (*domain.UsuarioDetail, error)
	RegistrarUsuario(ctx context.Context, usuarioRequest *domain.UsuarioRequest) (*domain.UsuarioDetail, error)
	ListarUsuarios(ctx context.Context) (*[]domain.UsuarioInfo, error)
	RestablecerPassword(ctx context.Context, usuarioId *int) (*domain.UsuarioDetail, error)
}

type UsuarioHandler interface {
	ObtenerUsuarioDetalle(c *fiber.Ctx) error
	RegistrarUsuario(c *fiber.Ctx) error
	HabilitarUsuarioById(c *fiber.Ctx) error
	DeshabilitarUsuarioById(c *fiber.Ctx) error
	ModificarUsuario(c *fiber.Ctx) error
	ListarUsuarios(c *fiber.Ctx) error
	ObtenerUsuarioActual(c *fiber.Ctx) error
	RestablecerPassword(c *fiber.Ctx) error
}
