package handler

import (
	"errors"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"
)

type AuthHandler struct {
	authService port.AuthService
}

func (a *AuthHandler) RefreshOrVerify(c *fiber.Ctx) error {

	claimsRefreshToken, err := util.Token.VerifyToken(c.Cookies("refresh-token"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(util.NewMessage("Usuario no autorizado"))
	}

	username, ok := claimsRefreshToken["username"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(util.NewMessage("Usuario no autorizado"))
	}
	userId, ok := claimsRefreshToken["userId"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(util.NewMessage("Usuario no autorizado"))
	}
	now := time.Now().UTC()
	expAccess, expRefresh := now.Add(1*time.Hour), now.Add(7*24*time.Hour)
	// Generar token
	accessToken, err := util.Token.CreateToken(jwt.MapClaims{
		"userId":     userId,
		"username":   username,
		"expiration": expAccess.Unix(),
		"type":       "access-token-adm",
	})

	// Generar token
	refreshToken, err := util.Token.CreateToken(jwt.MapClaims{
		"userId":     userId,
		"username":   username,
		"expiration": expRefresh.Unix(),
		"type":       "refresh-token-adm",
	})

	util.SetCookie(c, "access-token", accessToken, expAccess.Sub(now), true, false, now)
	util.SetCookie(c, "exp-access-token", fmt.Sprintf("%d", expAccess.Unix()), expAccess.Sub(now), false, false, now)

	expFloat, ok := claimsRefreshToken["expiration"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(util.NewMessage("Usuario no autorizado"))
	}
	expRefreshCurrent := time.Unix(int64(expFloat), 0)
	// REFRESH TOKEN: Renovar si queda menos de 6 horas
	if time.Until(expRefreshCurrent) < 6*time.Hour {
		util.SetCookie(c, "refresh-token", refreshToken, expRefresh.Sub(now), true, false, now)
		util.SetCookie(c, "exp-refresh-token", fmt.Sprintf("%d", expRefresh.Unix()), expRefresh.Sub(now), false, false, now)
	}

	return c.JSON(util.NewMessage("Sesión activa"))
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
