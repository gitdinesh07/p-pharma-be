package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ppharma/backend/pkg/api"
)

type InternalHandler struct{}

func NewInternalHandler() *InternalHandler { return &InternalHandler{} }

func (h *InternalHandler) SyncInventory(c *gin.Context) {
	c.JSON(http.StatusAccepted, api.APIResponse[map[string]string]{
		Success: true,
		Data:    &map[string]string{"status": "inventory sync accepted"},
	})
}

func (h *InternalHandler) BulkUpsertProducts(c *gin.Context) {
	c.JSON(http.StatusAccepted, api.APIResponse[map[string]string]{
		Success: true,
		Data:    &map[string]string{"status": "bulk upsert accepted"},
	})
}
