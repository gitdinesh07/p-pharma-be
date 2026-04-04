package handlers

import (
	"net/http"

	"ppharma/backend/internal/domain/common"
	"ppharma/backend/internal/domain/customer"
	"ppharma/backend/internal/http/middleware"
	"ppharma/backend/pkg/api"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	service *customer.Service
}

func NewCustomerHandler(service *customer.Service) *CustomerHandler {
	return &CustomerHandler{service: service}
}

type CreateCustomerRequest struct {
	Name     string              `json:"name" binding:"required"`
	Email    string              `json:"email"`
	Mobile   string              `json:"mobile"`
	Password string              `json:"password" binding:"required"`
	PhotoURL string              `json:"photo_url"`
	Gender   customer.GenderEnum `json:"gender"`
	Age      int                 `json:"age"`
	Address  customer.Address    `json:"address"`
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}

	var signupSource customer.SignupSource
	if val, ok := c.Get(middleware.ClientInfoKey); ok {
		info, _ := val.(common.ClientAppInfo)
		signupSource.Source = info.Source
		signupSource.ClientAppInfo.AppVersion = info.AppVersion
		signupSource.ClientAppInfo.DeviceType = info.DeviceType
		signupSource.ClientAppInfo.DeviceID = info.DeviceID
		signupSource.ClientAppInfo.OSVersion = ""
		signupSource.ClientAppInfo.DeviceModel = ""
	}

	cust := customer.Customer{
		Name:         req.Name,
		Email:        req.Email,
		Mobile:       req.Mobile,
		Password:     req.Password,
		PhotoURL:     req.PhotoURL,
		Gender:       req.Gender,
		Age:          req.Age,
		SignupSource: signupSource,
	}

	if req.Address.ValidateAddress() == nil {
		cust.Address = append(cust.Address, req.Address)
	}
	if err := h.service.CreateCustomer(&cust); err != nil {
		c.JSON(http.StatusInternalServerError, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "CREATE_FAILED", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, api.APIResponse[customer.Customer]{Success: true, Data: &cust})
}

func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id := c.Param("id")
	cust, err := h.service.GetCustomer(id)
	if err != nil {
		if err == customer.ErrCustomerNotFound {
			c.JSON(http.StatusNotFound, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "NOT_FOUND", Message: err.Error()}})
			return
		}
		c.JSON(http.StatusInternalServerError, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "FETCH_FAILED", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, api.APIResponse[customer.Customer]{Success: true, Data: cust})
}

func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	var cust customer.Customer
	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}
	cust.ID = id
	if err := h.service.UpdateCustomer(&cust); err != nil {
		c.JSON(http.StatusInternalServerError, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "UPDATE_FAILED", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, api.APIResponse[customer.Customer]{Success: true, Data: &cust})
}

type VerifyOTPRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	OTP        string `json:"otp" binding:"required"`
}

// func (h *CustomerHandler) Verify(c *gin.Context) {
// 	var req VerifyOTPRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "BAD_REQUEST", Message: err.Error()}})
// 		return
// 	}

// 	if err := h.service.VerifyOTP(req.Identifier, req.OTP); err != nil {
// 		c.JSON(http.StatusUnauthorized, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "VERIFY_FAILED", Message: err.Error()}})
// 		return
// 	}

// 	c.JSON(http.StatusOK, api.APIResponse[any]{Success: true, Message: "OTP verified successfully"})
// }
