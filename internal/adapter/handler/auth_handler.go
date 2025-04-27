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
	timeNow := time.Now().UTC()

	refreshToken := c.Cookies("refresh-token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(util.NewMessage("No se encontró el refresh token"))
	}

	claims, err := util.Token.VerifyToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(util.NewMessage("Refresh token inválido o expirado"))
	}

	username, ok := claims["username"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(util.NewMessage("Token inválido"))
	}

	// ACCESS TOKEN
	accessExp := timeNow.Add(15 * time.Minute)
	accessClaims := jwt.MapClaims{
		"username":   username,
		"expiration": accessExp.Unix(),
		"type":       "access",
	}

	newAccessToken, err := util.Token.CreateToken(accessClaims)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(util.NewMessage("No se pudo generar access token"))
	}

	util.SetCookie(c, "access-token", newAccessToken, 15*time.Minute, true, false, timeNow)
	util.SetCookie(c, "exp-access-token", fmt.Sprintf("%d", accessExp.Unix()), 15*time.Minute, false, false, timeNow)

	// REFRESH TOKEN: renovar si queda < 24h
	expFloat, ok := claims["expiration"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(util.NewMessage("Token inválido"))
	}
	refreshExp := time.Unix(int64(expFloat), 0)

	if time.Until(refreshExp) < 24*time.Hour {
		newRefreshExp := timeNow.Add(7 * 24 * time.Hour)
		newRefreshClaims := jwt.MapClaims{
			"username":   username,
			"expiration": newRefreshExp.Unix(),
			"type":       "refresh",
		}

		newRefreshToken, err := util.Token.CreateToken(newRefreshClaims)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(util.NewMessage("No se pudo generar nuevo refresh token"))
		}

		util.SetCookie(c, "refresh-token", newRefreshToken, 7*24*time.Hour, true, false, timeNow)
		util.SetCookie(c, "exp-refresh-token", fmt.Sprintf("%d", newRefreshExp.Unix()), 7*24*time.Hour, false, false, timeNow)
	} else {
		duration := time.Until(refreshExp)
		util.SetCookie(c, "refresh-token", refreshToken, duration, true, false, timeNow)
		util.SetCookie(c, "exp-refresh-token", fmt.Sprintf("%d", refreshExp.Unix()), duration, false, false, timeNow)
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
	tokenResponse, err := a.authService.ObtenerToken(ctx, &credentials)
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
	util.SetCookie(c, "refresh-token", tokenResponse.TokenRefresh, 7*24*time.Hour, true, false, now)
	util.SetCookie(c, "exp-refresh-token", fmt.Sprintf("%d", now.Add(7*24*time.Hour).Unix()), 7*24*time.Hour, false, false, now)

	// Access token: 15 minutos
	util.SetCookie(c, "access-token", tokenResponse.TokenAccess, 15*time.Minute, true, false, now)
	util.SetCookie(c, "exp-access-token", fmt.Sprintf("%d", now.Add(15*time.Minute).Unix()), 15*time.Minute, false, false, now)

	return c.JSON(util.NewMessage("Usuario autenticado"))
}

func NewAuthHandler(authService port.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

var _ port.AuthHandler = (*AuthHandler)(nil)
