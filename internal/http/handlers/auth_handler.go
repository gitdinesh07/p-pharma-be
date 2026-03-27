package handlers

import (
	"net/http"
	"time"

	"ppharma/backend/internal/domain/common"
	"ppharma/backend/internal/domain/customer"
	"ppharma/backend/internal/domain/user"
	"ppharma/backend/pkg/api"
	"ppharma/backend/support-pkg/auth/jwt"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	customerRepo customer.Repository
	userRepo     user.Repository
	jwtProvider  *jwt.Provider
}

func NewAuthHandler(customerRepo customer.Repository, userRepo user.Repository, jwtProvider *jwt.Provider) *AuthHandler {
	return &AuthHandler{
		customerRepo: customerRepo,
		userRepo:     userRepo,
		jwtProvider:  jwtProvider,
	}
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"` // Mobile or Email
	Password   string `json:"password" binding:"required"`
}

func isTestUser(identifier, password string) bool {
	return (identifier == "test@gmail.com" || identifier == "911122334455") && password == "test"
}

func (h *AuthHandler) CustomerLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}

	// MOCK VALIDATION: Hardcoded mock verification
	if !isTestUser(req.Identifier, req.Password) {
		c.JSON(http.StatusUnauthorized, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "UNAUTHORIZED", Message: "invalid credentials"}})
		return
	}

	token, err := h.jwtProvider.GenerateToken(&common.Principal{ID: "mock-customer-id", Role: "customer"}, 30*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "INTERNAL_ERROR", Message: "failed to generate token"}})
		return
	}

	c.JSON(http.StatusOK, api.APIResponse[map[string]string]{
		Success: true,
		Data:    &map[string]string{"access_token": token},
	})
}

func (h *AuthHandler) UserLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}

	// MOCK VALIDATION: Hardcoded mock verification
	if !isTestUser(req.Identifier, req.Password) {
		c.JSON(http.StatusUnauthorized, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "UNAUTHORIZED", Message: "invalid credentials"}})
		return
	}

	token, err := h.jwtProvider.GenerateToken(&common.Principal{ID: "mock-user-id", Role: "admin"}, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "INTERNAL_ERROR", Message: "failed to generate token"}})
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
