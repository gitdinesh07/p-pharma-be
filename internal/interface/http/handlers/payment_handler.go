package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ppharma/backend/pkg/api"
)

type PaymentHandler struct{}

func NewPaymentHandler() *PaymentHandler { return &PaymentHandler{} }

func (h *PaymentHandler) ListOrderPayments(c *gin.Context) {
	c.JSON(http.StatusOK, api.APIResponse[[]map[string]any]{
		Success: true,
		Data:    &[]map[string]any{{"payment_id": "pay_1", "order_id": c.Param("orderId"), "status": "pending", "method": "cod", "amount": 2100}},
	})
}

func (h *PaymentHandler) UpdatePaymentStatus(c *gin.Context) {
	c.JSON(http.StatusOK, api.APIResponse[map[string]string]{Success: true, Data: &map[string]string{"status": "payment status updated"}})
}
