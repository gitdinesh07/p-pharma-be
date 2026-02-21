package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ppharma/backend/pkg/api"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler { return &ProductHandler{} }

func (h *ProductHandler) ListProducts(c *gin.Context) {
	c.JSON(http.StatusOK, api.APIResponse[[]map[string]any]{
		Success: true,
		Data: &[]map[string]any{
			{"id": "prod_1", "name": "Med A", "sku": "MED-A", "inventory_count": 10},
			{"id": "prod_2", "name": "Med B", "sku": "MED-B", "inventory_count": 20},
		},
	})
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	c.JSON(http.StatusOK, api.APIResponse[map[string]any]{
		Success: true,
		Data:    &map[string]any{"id": c.Param("productId"), "name": "Med A", "sku": "MED-A", "inventory_count": 10},
	})
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	c.JSON(http.StatusCreated, api.APIResponse[map[string]string]{Success: true, Data: &map[string]string{"status": "product created"}})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	c.JSON(http.StatusOK, api.APIResponse[map[string]string]{Success: true, Data: &map[string]string{"status": "product updated"}})
}

func (h *ProductHandler) UpdateInventory(c *gin.Context) {
	c.JSON(http.StatusOK, api.APIResponse[map[string]string]{Success: true, Data: &map[string]string{"status": "inventory updated"}})
}
