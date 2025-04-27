package port

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"github.com/gofiber/fiber/v2"
)

type UsuarioRepository interface {
	ModificarEstadoUsuario(ctx context.Context, usuarioId *int) error
	ModificarUsuario(ctx context.Context, usuarioId *int, usuarioRequest *domain.UsuarioRequest) error
	ObtenerUsuario(ctx context.Context, username *string) (*domain.Usuario, error)
	ObtenerUsuarioDetalle(ctx context.Context, usuarioId *int) (*domain.UsuarioDetalle, error)
	RegistrarUsuario(ctx context.Context, usuarioRequest *domain.UsuarioRequest) (*domain.UsuarioDetalle, error)
	ListarUsuarios(ctx context.Context) (*[]domain.UsuarioInfo, error)
}

type UsuarioService interface {
	ModificarEstadoUsuario(ctx context.Context, usuarioId *int) error
	ModificarUsuario(ctx context.Context, usuarioId *int, usuarioRequest *domain.UsuarioRequest) error
	ObtenerUsuarioDetalle(ctx context.Context, usuarioId *int) (*domain.UsuarioDetalle, error)
	RegistrarUsuario(ctx context.Context, usuarioRequest *domain.UsuarioRequest) (*domain.UsuarioDetalle, error)
	ListarUsuarios(ctx context.Context) (*[]domain.UsuarioInfo, error)
}

type UsuarioHandler interface {
	ObtenerUsuarioDetalle(c *fiber.Ctx) error
	RegistrarUsuario(c *fiber.Ctx) error
	ModificarEstadoUsuario(c *fiber.Ctx) error
	ModificarUsuario(c *fiber.Ctx) error
	ListarUsuarios(c *fiber.Ctx) error
}
