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
	return &ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: "Ha ocurrido un error interno en el servidor. Por favor, inténtelo más tarde.",
	}
}

func NewStatusServiceUnavailableErrorGeneric() *ErrorResponse {
	return &ErrorResponse{
		Code:    http.StatusServiceUnavailable,
		Message: "Servicio no disponible, inténtelo más tarde.",
	}
}
func NewStatusUnauthorizedError(message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    http.StatusUnauthorized,
		Message: message,
	}

}
func NewInternalServerError(message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}
func NewBadRequestError(message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewConflictError(message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    http.StatusConflict,
		Message: message,
	}
}

func NewNotFoundError(message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

//func NewErrorResponse(code int, message string) *ErrorResponse {
//	return &ErrorResponse{
//		Code:    code,
//		Message: message,
//	}
//}

var _ error = (*ErrorResponse)(nil)
