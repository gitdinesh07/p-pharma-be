package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"ppharma/backend/internal/domain/common"
	"ppharma/backend/pkg/api"
)

type InternalHandler struct {
	queue common.Queue
	topic string
}

func NewInternalHandler(queue common.Queue, topic string) *InternalHandler {
	return &InternalHandler{queue: queue, topic: topic}
}

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

func (h *InternalHandler) PublishQueueMessage(c *gin.Context) {
	payload, _ := json.Marshal(map[string]any{
		"trigger": "api_internal_route",
		"at":      time.Now().UTC().Format(time.RFC3339),
	})
	if err := h.queue.Publish(c.Request.Context(), h.topic, payload); err != nil {
		c.JSON(http.StatusInternalServerError, api.APIResponse[any]{
			Success: false,
			Error:   &api.APIError{Code: "QUEUE_PUBLISH_FAILED", Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusAccepted, api.APIResponse[map[string]string]{
		Success: true,
		Data:    &map[string]string{"status": "message published"},
	})
}
