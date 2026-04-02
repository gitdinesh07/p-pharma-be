package handlers

import (
	"net/http"

	"ppharma/backend/internal/domain/customer"
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

	cust := customer.Customer{
		Name:     req.Name,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Password: req.Password,
		PhotoURL: req.PhotoURL,
		Gender:   req.Gender,
		Age:      req.Age,
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
