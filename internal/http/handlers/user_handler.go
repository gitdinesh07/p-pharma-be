package handlers

import (
	"net/http"

	"ppharma/backend/internal/domain/user"
	"ppharma/backend/pkg/api"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *user.Service
}

func NewUserHandler(userService *user.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type CreateUserRequest struct {
	Name     string    `json:"name" binding:"required"`
	Email    string    `json:"email"`
	Mobile   string    `json:"mobile"`
	Password string    `json:"password" binding:"required"`
	PhotoURL string    `json:"photo_url"`
	Role     user.Role `json:"role" binding:"required"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}

	u := &user.User{
		Name:     req.Name,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Password: req.Password,
		PhotoURL: req.PhotoURL,
		Role:     req.Role,
	}

	if err := h.userService.CreateUser(u); err != nil {
		c.JSON(http.StatusConflict, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "CONFLICT", Message: err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, api.APIResponse[user.User]{
		Success: true,
		Data:    u,
	})
}

type UpdateUserRequest struct {
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Mobile   string    `json:"mobile"`
	PhotoURL string    `json:"photo_url"`
	Role     user.Role `json:"role"`
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "BAD_REQUEST", Message: "id is required"}})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}

	u, err := h.userService.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "NOT_FOUND", Message: err.Error()}})
		return
	}

	if req.Name != "" {
		u.Name = req.Name
	}
	if req.Email != "" {
		u.Email = req.Email
	}
	if req.Mobile != "" {
		u.Mobile = req.Mobile
	}
	if req.PhotoURL != "" {
		u.PhotoURL = req.PhotoURL
	}
	if req.Role != "" {
		u.Role = req.Role
	}

	if err := h.userService.UpdateUser(u); err != nil {
		c.JSON(http.StatusConflict, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "CONFLICT", Message: err.Error()}})
		return
	}

	c.JSON(http.StatusOK, api.APIResponse[user.User]{
		Success: true,
		Data:    u,
	})
}

type VerifyUserOTPRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	OTP        string `json:"otp" binding:"required"`
}

func (h *UserHandler) VerifyOTP(c *gin.Context) {
	var req VerifyUserOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}

	if err := h.userService.VerifyOTP(req.Identifier, req.OTP); err != nil {
		c.JSON(http.StatusUnauthorized, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "VERIFY_FAILED", Message: err.Error()}})
		return
	}

	c.JSON(http.StatusOK, api.APIResponse[any]{Success: true, Message: "OTP verified successfully"})
}
