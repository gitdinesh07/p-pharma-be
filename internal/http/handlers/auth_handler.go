package handlers

import (
	"errors"
	"net/http"

	"ppharma/backend/internal/service"
	"ppharma/backend/pkg/api"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"` // Mobile or Email
	Password   string `json:"password" binding:"required"`
}

func (h *AuthHandler) CustomerLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}

	token, err := h.authService.CustomerLogin(req.Identifier, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "UNAUTHORIZED", Message: err.Error()}})
		} else {
			c.JSON(http.StatusInternalServerError, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "INTERNAL_ERROR", Message: err.Error()}})
		}
		return
	}

	c.JSON(http.StatusOK, api.APIResponse[map[string]string]{
		Success: true,
		Data:    &map[string]string{"access_token": token},
	})
}

type ResetPasswordRequest struct {
	Identifier  string `json:"identifier" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (h *AuthHandler) ResetCustomerPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}

	err := h.authService.ResetCustomerPassword(req.Identifier, req.NewPassword)
	if err != nil {
		if errors.Is(err, service.ErrCustomerNotFound) {
			c.JSON(http.StatusNotFound, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "NOT_FOUND", Message: err.Error()}})
		} else {
			c.JSON(http.StatusInternalServerError, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "INTERNAL_ERROR", Message: err.Error()}})
		}
		return
	}

	c.JSON(http.StatusOK, api.APIResponse[map[string]string]{
		Success: true,
		Data:    &map[string]string{"status": "password reset successfully"},
	})
}

func (h *AuthHandler) UserLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}

	token, err := h.authService.UserLogin(req.Identifier, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "UNAUTHORIZED", Message: err.Error()}})
		} else {
			c.JSON(http.StatusInternalServerError, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "INTERNAL_ERROR", Message: err.Error()}})
		}
		return
	}

	c.JSON(http.StatusOK, api.APIResponse[map[string]string]{
		Success: true,
		Data:    &map[string]string{"access_token": token},
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	c.JSON(http.StatusOK, api.APIResponse[map[string]string]{
		Success: true,
		Data:    &map[string]string{"access_token": "mock-access-token-refreshed"},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, api.APIResponse[map[string]string]{Success: true, Data: &map[string]string{"status": "logged out"}})
}

func (h *AuthHandler) ListSessions(c *gin.Context) {
	c.JSON(http.StatusOK, api.APIResponse[[]map[string]string]{Success: true, Data: &[]map[string]string{{"session_id": "session-1", "device": "web"}}})
}

func (h *AuthHandler) RevokeSession(c *gin.Context) {
	c.JSON(http.StatusOK, api.APIResponse[map[string]string]{Success: true, Data: &map[string]string{"status": "session revoked"}})
}
