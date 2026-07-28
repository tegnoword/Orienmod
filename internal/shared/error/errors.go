package errors

import (
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

type AppError struct {
	Message string
	Code    int
	Err     error
	Details string
}

func NewBadRequest(msg string) *AppError {
	return &AppError{
		Message: msg,
		Code:    http.StatusBadRequest,
	}
}

func NewBadRequestWithDetails(msg, details string) *AppError {
	return &AppError{
		Message: msg,
		Code:    http.StatusBadRequest,
		Details: details,
	}
}

func NewUnauthorized(msg string) *AppError {
	return &AppError{
		Message: msg,
		Code:    http.StatusUnauthorized,
	}
}

func NewForbidden(msg string) *AppError {
	return &AppError{
		Message: msg,
		Code:    http.StatusForbidden,
	}
}

func NewNotFound(msg string) *AppError {
	return &AppError{
		Message: msg,
		Code:    http.StatusNotFound,
	}
}

func NewConflict(msg string) *AppError {
	return &AppError{
		Message: msg,
		Code:    http.StatusConflict,
	}
}

func NewInternalError(err error) *AppError {
	return &AppError{
		Message: "Error interno del servidor",
		Code:    http.StatusInternalServerError,
		Err:     err,
	}
}

func NewInternalErrorWithMsg(msg string, err error) *AppError {
	return &AppError{
		Message: msg,
		Code:    http.StatusInternalServerError,
		Err:     err,
	}
}

func NewValidationError(field, msg string) *AppError {
	return &AppError{
		Message: fmt.Sprintf("Error de validación en '%s': %s", field, msg),
		Code:    http.StatusBadRequest,
		Details: fmt.Sprintf("Campo: %s", field),
	}
}

func NewGoogleAPIError(err error) *AppError {
	return &AppError{
		Message: "Error al comunicarse con Google Classroom",
		Code:    http.StatusBadGateway,
		Err:     err,
	}
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	if e.Details != "" {
		return fmt.Sprintf("%s - %s", e.Message, e.Details)
	}
	return e.Message
}

func (e *AppError) ToResponse() ErrorResponse {
	resp := ErrorResponse{
		Error: e.Message,
		Code:  e.Code,
	}
	if e.Details != "" {
		resp.Details = e.Details
	}
	if e.Err != nil && e.Code >= 500 {
		resp.Details = e.Err.Error()
	}
	return resp
}

func (e *AppError) IsInternal() bool {
	return e.Code >= 500
}

func WrapError(err error, msg string) *AppError {
	return &AppError{
		Message: msg,
		Code:    http.StatusInternalServerError,
		Err:     err,
	}
}

func HandleError(err error) (int, interface{}) {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code, appErr.ToResponse()
	}

	return http.StatusInternalServerError, ErrorResponse{
		Error:   "Error interno del servidor",
		Code:    http.StatusInternalServerError,
		Details: err.Error(),
	}
}
