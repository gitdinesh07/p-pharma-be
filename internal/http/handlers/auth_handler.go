package handlers

import (
	"errors"
	"net/http"
	"ppharma/backend/internal/domain/auth"
	"ppharma/backend/pkg/api"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *auth.AuthService
}

func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type LoginRequest struct {
	Identifier     string `json:"identifier" binding:"required"` // Mobile or Email
	Password       string `json:"password"`
	Otp            string `json:"otp"`
	IsAdminUserReq bool   `json:"is_admin_user_req" binding:"required"`
}

type ResetPasswordRequest struct {
	Identifier  string `json:"identifier" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type VerifyCustomerOtpGenerateTokenRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	OTP        string `json:"otp" binding:"required"`
}

// func (h *AuthHandler) CustomerRegister(c *gin.Context) {
// 	var req LoginRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "BAD_REQUEST", Message: err.Error()}})
// 		return
// 	}

// 	token, err := h.authService.CustomerRegister(req.Identifier, req.Password, req.Otp)
// 	if err != nil {
// 		if errors.Is(err, auth.ErrInvalidCredentials) {
// 			c.JSON(http.StatusUnauthorized, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "UNAUTHORIZED", Message: err.Error()}})
// 		} else {
// 			c.JSON(http.StatusInternalServerError, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "INTERNAL_ERROR", Message: err.Error()}})
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, api.APIResponse[map[string]string]{
// 		Success: true,
// 		Data:    &map[string]string{"access_token": token},
// 	})
// }

// func (h *AuthHandler) CustomerSendOTP(c *gin.Context) {
// 	var req LoginRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "BAD_REQUEST", Message: err.Error()}})
// 		return
// 	}

// 	err := h.authService.CustomerSendOTP(req.Identifier)
// 	if err != nil {
// 		if errors.Is(err, auth.ErrInvalidCredentials) {
// 			c.JSON(http.StatusUnauthorized, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "UNAUTHORIZED", Message: err.Error()}})
// 		} else {
// 			c.JSON(http.StatusInternalServerError, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "INTERNAL_ERROR", Message: err.Error()}})
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, api.APIResponse[map[string]string]{
// 		Success: true,
// 		Data:    &map[string]string{"status": "otp sent successfully"},
// 	})
// }

func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}

	var token string
	var err error
	//admin user login handler
	if req.IsAdminUserReq {
		token, err = h.authService.UserLogin(req.Identifier, req.Password)
	}

	//customer login handler
	if !req.IsAdminUserReq {
		token, err = h.authService.CustomerLogin(req.Identifier, req.Password, req.Otp)
	}

	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
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

func (h *AuthHandler) ResetCustomerPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}

	err := h.authService.ResetCustomerPassword(req.Identifier, req.NewPassword)
	if err != nil {
		if errors.Is(err, auth.ErrCustomerNotFound) {
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
		if errors.Is(err, auth.ErrInvalidCredentials) {
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

func (h *AuthHandler) GenerateAndSendOtp(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}

	err := h.authService.SendVerificationOtp(req.Identifier)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "UNAUTHORIZED", Message: err.Error()}})
		} else {
			c.JSON(http.StatusInternalServerError, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "INTERNAL_ERROR", Message: err.Error()}})
		}
		return
	}

	c.JSON(http.StatusOK, api.APIResponse[map[string]string]{
		Success: true,
		Data:    &map[string]string{"status": "otp sent successfully"},
	})
}
