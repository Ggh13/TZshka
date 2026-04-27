package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	historyService "TZshka/internal/history/service"
)

type Handler struct {
	service *historyService.Service
	log     *zap.Logger
}

func New(service *historyService.Service, log *zap.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

func (h *Handler) GetHistory(c *gin.Context) {
	sessionIDStr := c.Param("id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		h.log.Error("invalid session id", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	history, err := h.service.GetHistory(c.Request.Context(), sessionID)
	if err != nil {
		h.log.Error("failed to get history", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"corrections": history})
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.GET("/history/:id", h.GetHistory)
}
