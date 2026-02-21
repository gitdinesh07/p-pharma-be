package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ppharma/backend/pkg/api"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler { return &AuthHandler{} }

func (h *AuthHandler) Login(c *gin.Context) {
	c.JSON(http.StatusOK, api.APIResponse[map[string]string]{
		Success: true,
		Data: &map[string]string{
			"access_token":  "mock-access-token",
			"refresh_token": "mock-refresh-token",
		},
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
