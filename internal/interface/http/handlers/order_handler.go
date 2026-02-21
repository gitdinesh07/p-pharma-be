package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ppharma/backend/internal/domain/auth"
	"ppharma/backend/internal/domain/common"
	"ppharma/backend/internal/domain/order"
	"ppharma/backend/pkg/api"
)

type OrderHandler struct {
	service *order.Service
}

func NewOrderHandler(service *order.Service) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) GetCustomerOrder(c *gin.Context) {
	v, _ := c.Get(string(auth.ContextPrincipalKey))
	principal, _ := v.(*common.Principal)
	orderID := c.Param("orderId")
	ord, err := h.service.GetOrderForCustomer(orderID, principal.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "NOT_FOUND", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, api.APIResponse[order.Order]{Success: true, Data: ord})
}

type updateItemStatusRequest struct {
	Status string `json:"status" binding:"required"`
	Reason string `json:"reason"`
}

func (h *OrderHandler) UpdateOrderItemStatus(c *gin.Context) {
	var req updateItemStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}
	status := order.ItemStatus(req.Status)
	orderID := c.Param("orderId")
	itemID := c.Param("itemId")
	changedBy := "admin"
	if v, ok := c.Get(string(auth.ContextPrincipalKey)); ok {
		if p, ok := v.(*common.Principal); ok {
			changedBy = p.ID
		}
	}
	if v, ok := c.Get(string(auth.ContextAPIKeyKey)); ok {
		if p, ok := v.(*common.APIKeyPrincipal); ok {
			changedBy = "internal:" + p.KeyID
		}
	}

	ord, err := h.service.UpdateItemStatus(orderID, itemID, status, req.Reason, changedBy)
	if err != nil {
		code := http.StatusBadRequest
		if err == order.ErrOrderNotFound || err == order.ErrItemNotFound {
			code = http.StatusNotFound
		}
		c.JSON(code, api.APIResponse[any]{Success: false, Error: &api.APIError{Code: "ORDER_UPDATE_FAILED", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, api.APIResponse[order.Order]{Success: true, Data: ord})
}
