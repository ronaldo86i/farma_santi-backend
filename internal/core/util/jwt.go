package util

import (
	"farma-santi_backend/internal/core/domain/datatype"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

type token struct{}

var Token token

// secretKeyJwtAdmin es la clave secreta que permite generar y extraer el token
var secretKeyJwtAdmin = []byte(os.Getenv("SECRET_KEY"))

// CreateToken genera un token JWT válido
func (token) CreateToken(data jwt.MapClaims) (string, error) {
	// Crear un nuevo token con la función de firma HS256 y los datos del usuario.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		data,
	)

	// Firmar el token con la clave secreta.
	tokenString, err := token.SignedString(secretKeyJwtAdmin)
	if err != nil {
		return "", fmt.Errorf("error al firmar el token: %w", err) // Retornar error si no se puede firmar el token.
	}

	// Retornar el token firmado.
	return tokenString, nil
}

// VerifyToken verifica si un token JWT es válido y retorna los claims.
func (token) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	var message = "Token no válido"
	// Parsear el token y validar la firma con la clave secreta.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verificar que la firma sea HS256.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			message = "Método de firma inválido"
			return nil, &datatype.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: message,
			}
		}
		return secretKeyJwtAdmin, nil
	})
	// Si hay un error en el parseo, retornar el error.
	if err != nil {
		message = "Error al verificar el token"
		return nil, &datatype.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: message,
		}
	}
	// Verificar si el token es válido.
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Verificar la expiración del token.
		if exp, ok := claims["expiration"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				message = "El token ha expirado"
				return nil, &datatype.ErrorResponse{
					Code:    http.StatusUnauthorized,
					Message: message,
				}
			}
		} else {
			message = "Error al verificar el token"
			return nil, &datatype.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: message,
			}
		}
		// Retornar los claims del token.
		return claims, nil
	}
	return nil, &datatype.ErrorResponse{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}
