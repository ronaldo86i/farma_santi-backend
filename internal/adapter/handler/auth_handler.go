package handler

import (
	"context"
	"errors"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/service"
	"farma-santi_backend/internal/core/util"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	authService port.AuthService
}

func (a *AuthHandler) RegisterWithEmail(c *fiber.Ctx) error {
	var credentials domain.FirebaseLogin
	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("Datos no válidos"))
	}

	fb := service.GetFirebaseClient()
	// Verificar token de Firebase
	token, err := fb.AuthClient.VerifyIDToken(context.Background(), credentials.Token)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(util.NewMessage("Token inválido"))
	}

	// 🔎 Verificar si el email ya fue confirmado
	if verified, ok := token.Claims["email_verified"].(bool); !ok || !verified {
		return c.Status(http.StatusUnauthorized).JSON(util.NewMessage("Correo no verificado"))
	}

	// Crear JWT propio
	expAccess := time.Now().UTC().Add(24 * time.Hour)
	accessToken, err := util.Token.CreateToken(jwt.MapClaims{
		"userId":     token.UID,
		"email":      token.Claims["email"],
		"type":       "access-token-public",
		"expiration": expAccess.Unix(),
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage("Error generando token"))
	}

	return c.JSON(fiber.Map{"token": accessToken})
}

func (a *AuthHandler) LoginWithGoogle(c *fiber.Ctx) error {
	var credentials domain.FirebaseLogin
	if err := c.BodyParser(&credentials); err != nil {
		return c.JSON(util.NewMessage("Datos no válidos"))
	}
	fb := service.GetFirebaseClient()
	token, err := fb.AuthClient.VerifyIDToken(context.Background(), credentials.Token)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage(err.Error()))
	}
	// También puedes obtener el correo del token
	email := token.Claims["email"].(string)

	expAccess := time.Now().UTC().Add(60 * 24 * time.Hour)
	accessToken, err := util.Token.CreateToken(jwt.MapClaims{
		"userId":     token.UID,
		"email":      email,
		"expiration": expAccess.Unix(),
		"type":       "access-token-public",
	})
	return c.JSON(fiber.Map{"token": accessToken})
}

func (a *AuthHandler) LoginWithEmail(c *fiber.Ctx) error {
	var credentials domain.FirebaseLogin
	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(http.StatusBadRequest).JSON(util.NewMessage("Datos no válidos"))
	}

	fb := service.GetFirebaseClient()
	token, err := fb.AuthClient.VerifyIDToken(context.Background(), credentials.Token)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(util.NewMessage("Token inválido"))
	}

	// Verificar que el correo esté validado en Firebase
	if !token.Claims["email_verified"].(bool) {
		return c.Status(http.StatusForbidden).JSON(util.NewMessage("El correo no está verificado"))
	}

	// También puedes obtener el correo del token
	email := token.Claims["email"].(string)

	// Crear JWT propio
	expAccess := time.Now().UTC().Add(60 * 24 * time.Hour)
	accessToken, err := util.Token.CreateToken(jwt.MapClaims{
		"userId":     token.UID,
		"email":      email,
		"expiration": expAccess.Unix(),
		"type":       "access-token-public",
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage("No se pudo crear token"))
	}

	return c.JSON(fiber.Map{"token": accessToken})
}

func (a *AuthHandler) RefreshOrVerify(c *fiber.Ctx) error {
	now := time.Now().UTC()

	// Verificar refresh-token
	claims, err := util.Token.VerifyToken(c.Cookies("refresh-token"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"authenticated": false, "message": "Usuario no autorizado"})
	}

	username, ok := claims["username"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"authenticated": false})
	}
	userId, ok := claims["userId"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"authenticated": false})
	}

	// Expiraciones
	expAccess := now.Add(1 * time.Hour)
	expRefresh := now.Add(7 * 24 * time.Hour)

	// Crear tokens
	accessToken, _ := util.Token.CreateToken(jwt.MapClaims{
		"userId":     userId,
		"username":   username,
		"expiration": expAccess.Unix(),
		"type":       "access-token-adm",
	})
	refreshToken, _ := util.Token.CreateToken(jwt.MapClaims{
		"userId":     userId,
		"username":   username,
		"expiration": expRefresh.Unix(),
		"type":       "refresh-token-adm",
	})

	// Guardar cookies cross-domain
	util.SetCookie(c, "access-token", accessToken, expAccess.Sub(now), true, false, now)
	util.SetCookie(c, "exp-access-token", fmt.Sprintf("%d", expAccess.Unix()), expAccess.Sub(now), false, false, now)

	// Renovar refresh token si queda menos de 6 horas
	expFloat, _ := claims["expiration"].(float64)
	expRefreshCurrent := time.Unix(int64(expFloat), 0)
	if time.Until(expRefreshCurrent) < 6*time.Hour {
		util.SetCookie(c, "refresh-token", refreshToken, expRefresh.Sub(now), true, false, now)
		util.SetCookie(c, "exp-refresh-token", fmt.Sprintf("%d", expRefresh.Unix()), expRefresh.Sub(now), false, false, now)
	}

	// Respuesta explícita
	return c.JSON(fiber.Map{"authenticated": true, "message": "Sesión activa"})
}

func (a *AuthHandler) Logout(c *fiber.Ctx) error {
	util.DeleteCookie(c, "access-token", true)
	util.DeleteCookie(c, "refresh-token", true)
	util.DeleteCookie(c, "exp-access-token", false)
	util.DeleteCookie(c, "exp-refresh-token", false)

	return c.JSON(util.NewMessage("Sesión finalizada"))
}

func (a *AuthHandler) Login(c *fiber.Ctx) error {
	var credentials domain.LoginRequest
	if err := c.BodyParser(&credentials); err != nil {
		return c.JSON(util.NewMessage("Datos no válidos"))
	}

	ctx := c.UserContext()
	tokenResponse, err := a.authService.ObtenerTokenByCredencial(ctx, &credentials)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}

	now := time.Now().UTC()

	// Refresh token: 7 días
	util.SetCookie(c, "refresh-token", tokenResponse.RefreshToken, tokenResponse.ExpRefreshToken.Sub(now), true, false, now)
	util.SetCookie(c, "exp-refresh-token", fmt.Sprintf("%d", tokenResponse.ExpRefreshToken.Unix()), tokenResponse.ExpRefreshToken.Sub(now), false, false, now)

	// Access token: 15 minutos
	util.SetCookie(c, "access-token", tokenResponse.AccessToken, tokenResponse.ExpAccessToken.Sub(now), true, false, now)
	util.SetCookie(c, "exp-access-token", fmt.Sprintf("%d", tokenResponse.ExpAccessToken.Unix()), tokenResponse.ExpAccessToken.Sub(now), false, false, now)

	return c.JSON(util.NewMessage("Usuario autenticado"))
}

func NewAuthHandler(authService port.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

var _ port.AuthHandler = (*AuthHandler)(nil)
