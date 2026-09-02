package http

import (
	"net/http"

	"kopelko-dating-app-backend/internal/core/identity/domain"
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
	Email     string   `json:"email" validate:"required,email" example:"user@example.com"`
	Password  string   `json:"password" validate:"required,min=6" example:"secret123"`
	Name      string   `json:"name" validate:"required" example:"John Doe"`
	Age       int      `json:"age" validate:"required,gte=18" example:"25"`
	Gender    string   `json:"gender" validate:"required,oneof=male female other" example:"male"`
	Location  string   `json:"location" validate:"required" example:"Jakarta"`
	Interests []string `json:"interests" validate:"omitempty" example:"coding,music"`
	Photos    []string `json:"photos" validate:"omitempty" example:"https://picsum.photos/200"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required" example:"secret123"`
}

type AuthResponse struct {
	User  *domain.User `json:"user"`
	Token string       `json:"token"`
}

// Register godoc
// @Summary Register new user
// @Description Register a new user and create their initial profile atomically
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register payload"
// @Success 201 {object} pkgHttp.APIResponse{data=AuthResponse}
// @Failure 400 {object} pkgHttp.APIResponse
// @Router /api/register [post]
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

// Login godoc
// @Summary Login user
// @Description Authenticate user by email & password to obtain JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login payload"
// @Success 200 {object} pkgHttp.APIResponse{data=AuthResponse}
// @Failure 400 {object} pkgHttp.APIResponse
// @Failure 401 {object} pkgHttp.APIResponse
// @Router /api/login [post]
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
