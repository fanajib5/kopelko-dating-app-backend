package http

import (
	"net/http"

	"kopelko-dating-app-backend/internal/modules/identity/domain"
	pkgHttp "kopelko-dating-app-backend/internal/platform/http"

	"github.com/labstack/echo/v4"
)

type IdentityHandler struct {
	svc domain.IdentityService
}

func NewIdentityHandler(svc domain.IdentityService) *IdentityHandler {
	return &IdentityHandler{svc: svc}
}

type RegisterRequest struct {
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required,min=6"`
	Name      string   `json:"name" validate:"required"`
	Age       int      `json:"age" validate:"required,gte=18"`
	Gender    string   `json:"gender" validate:"required,oneof=male female other"`
	Location  string   `json:"location" validate:"required"`
	Interests []string `json:"interests" validate:"omitempty"`
	Photos    []string `json:"photos" validate:"omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User  *domain.User `json:"user"`
	Token string       `json:"token"`
}

func (h *IdentityHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return pkgHttp.BadRequest(c, "Invalid request body", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return pkgHttp.BadRequest(c, "Validation failed", err.Error())
	}

	user, token, err := h.svc.Register(
		c.Request().Context(),
		req.Email, req.Password, req.Name,
		req.Age, req.Gender, req.Location,
		req.Interests, req.Photos,
	)
	if err != nil {
		return pkgHttp.BadRequest(c, err.Error(), nil)
	}

	return pkgHttp.Success(c, http.StatusCreated, "User registered successfully", AuthResponse{
		User:  user,
		Token: token,
	})
}

func (h *IdentityHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return pkgHttp.BadRequest(c, "Invalid request body", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return pkgHttp.BadRequest(c, "Validation failed", err.Error())
	}

	user, token, err := h.svc.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return pkgHttp.Unauthorized(c, err.Error())
	}

	return pkgHttp.Success(c, http.StatusOK, "Login successful", AuthResponse{
		User:  user,
		Token: token,
	})
}

func (h *IdentityHandler) RegisterRoutes(api *echo.Group) {
	api.POST("/register", h.Register)
	api.POST("/login", h.Login)
}
