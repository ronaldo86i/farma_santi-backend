package datatype

import (
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Implementar la función `Error` para que cumpla con la interfaz `error`
func (e ErrorResponse) Error() string {
	return e.Message
}

// NewInternalServerError crea un objeto de error común para el servidor
func NewInternalServerError() *ErrorResponse {
	return &ErrorResponse{
		Code:    fiber.StatusInternalServerError,
		Message: "Ha ocurrido un error interno en el servidor. Por favor, inténtelo más tarde.",
	}
}

// NewStatusServiceUnavailableError crea un objeto de error de servicio no disponible para el servidor
func NewStatusServiceUnavailableError() *ErrorResponse {
	return &ErrorResponse{
		Code:    fiber.StatusServiceUnavailable,
		Message: "Servicio no disponible, inténtelo más tarde.",
	}
}

func NewErrorResponse(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}

var _ error = (*ErrorResponse)(nil)
