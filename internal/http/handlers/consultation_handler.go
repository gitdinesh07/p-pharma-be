package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ppharma/backend/pkg/api"
)

type ConsultationHandler struct{}

func NewConsultationHandler() *ConsultationHandler { return &ConsultationHandler{} }

func (h *ConsultationHandler) CreateConsultation(c *gin.Context) {
	c.JSON(http.StatusCreated, api.APIResponse[map[string]string]{
		Success: true,
		Data: &map[string]string{
			"consultation_id": "con_1",
			"status":          "scheduled",
			"meeting_url":     "https://meet.google.com/mock-link",
		},
	})
}

func (h *ConsultationHandler) GetConsultation(c *gin.Context) {
	c.JSON(http.StatusOK, api.APIResponse[map[string]any]{
		Success: true,
		Data: &map[string]any{
			"consultation_id": c.Param("consultationId"),
			"status":          "scheduled",
			"meeting_url":     "https://meet.google.com/mock-link",
		},
	})
}

func (h *ConsultationHandler) ListConsultations(c *gin.Context) {
	c.JSON(http.StatusOK, api.APIResponse[[]map[string]any]{
		Success: true,
		Data:    &[]map[string]any{{"consultation_id": "con_1", "status": "scheduled"}},
	})
}

func (h *ConsultationHandler) UpdateConsultationStatus(c *gin.Context) {
	c.JSON(http.StatusOK, api.APIResponse[map[string]string]{Success: true, Data: &map[string]string{"status": "consultation status updated"}})
}
