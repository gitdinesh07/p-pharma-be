package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ppharma/backend/internal/domain/customer"
	"ppharma/backend/pkg/api"
)

type CustomerHandler struct {
	service *customer.Service
}

func NewCustomerHandler(service *customer.Service) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var cust customer.Customer
	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "BAD_REQUEST", Message: err.Error()}})
		return
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
