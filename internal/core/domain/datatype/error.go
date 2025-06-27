package datatype

import (
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

func NewInternalServerErrorGeneric() *ErrorResponse {
	return NewErrorResponse(http.StatusInternalServerError, "Ha ocurrido un error interno en el servidor. Por favor, inténtelo más tarde.")
}

func NewStatusServiceUnavailableErrorGeneric() *ErrorResponse {
	return NewErrorResponse(http.StatusServiceUnavailable, "Servicio no disponible, inténtelo más tarde.")
}
func NewStatusUnauthorizedError(message string) *ErrorResponse {
	return NewErrorResponse(http.StatusUnauthorized, message)

}
func NewInternalServerError(message string) *ErrorResponse {
	return NewErrorResponse(http.StatusInternalServerError, message)
}
func NewBadRequestError(message string) *ErrorResponse {
	return NewErrorResponse(http.StatusBadRequest, message)
}

func NewConflictError(message string) *ErrorResponse {
	return NewErrorResponse(http.StatusConflict, message)
}

func NewNotFoundError(message string) *ErrorResponse {
	return NewErrorResponse(http.StatusNotFound, message)
}

// NewErrorResponse permite crear respuesta de error personalizado
func NewErrorResponse(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}

var _ error = (*ErrorResponse)(nil)
