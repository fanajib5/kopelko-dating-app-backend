package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func Success(c echo.Context, statusCode int, message string, data any) error {
	return c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c echo.Context, statusCode int, message string, err any) error {
	return c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}

func BadRequest(c echo.Context, message string, err any) error {
	return Error(c, http.StatusBadRequest, message, err)
}

func Unauthorized(c echo.Context, message string) error {
	return Error(c, http.StatusUnauthorized, message, nil)
}

func NotFound(c echo.Context, message string) error {
	return Error(c, http.StatusNotFound, message, nil)
}

func InternalServerError(c echo.Context, message string, err any) error {
	return Error(c, http.StatusInternalServerError, message, err)
}
