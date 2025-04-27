package service

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type AuthService struct {
	usuarioRepository port.UsuarioRepository
}

func (a AuthService) ObtenerToken(ctx context.Context, credentials *domain.LoginRequest) (*domain.TokenResponse, error) {
	usuario, err := a.usuarioRepository.ObtenerUsuario(ctx, &credentials.Username)
	if err != nil {
		return nil, err
	}
	// Comparar contraseña hashed y contraseña ingresada
	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(credentials.Password)); err != nil {
		return nil, &datatype.ErrorResponse{
			Code:    fiber.StatusUnauthorized,
			Message: "Usuario o contraseña incorrecta",
		}
	}
	expAccess, expRefresh := time.Now().UTC().Add(15*time.Minute).Unix(), time.Now().UTC().Add(7*24*time.Hour).Unix()
	// Generar token
	tokenAccess, err := util.Token.CreateToken(jwt.MapClaims{
		"username":   &credentials.Username,
		"expiration": expAccess,
		"type":       "access",
	})
	// Generar token
	tokenRefresh, err := util.Token.CreateToken(jwt.MapClaims{
		"username":   &credentials.Username,
		"expiration": expRefresh,
		"type":       "refresh",
	})
	if err != nil {
		return nil, &datatype.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Usuario o contraseña incorrecto",
		}
	}

	return &domain.TokenResponse{
		Message:      "Usuario autenticado",
		TokenAccess:  tokenAccess,
		TokenRefresh: tokenRefresh,
	}, nil
}

func NewAuthService(usuarioRepository port.UsuarioRepository) *AuthService {
	return &AuthService{usuarioRepository}
}

var _ port.AuthService = (*AuthService)(nil)
