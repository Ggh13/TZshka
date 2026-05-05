package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	historyModels "TZshka/internal/history/models"
	historyService "TZshka/internal/history/service"
	sessionRepo "TZshka/internal/session/repository"
	"TZshka/pkg/registr"
)

type Handler struct {
	service    *historyService.Service
	sessionRepo sessionRepo.Repository
	tokenSvc   *registr.TokenService
	log        *zap.Logger
}

func New(service *historyService.Service, sessionRepo sessionRepo.Repository, tokenSvc *registr.TokenService, log *zap.Logger) *Handler {
	return &Handler{
		service:    service,
		sessionRepo: sessionRepo,
		tokenSvc:   tokenSvc,
		log:        log,
	}
}

func (h *Handler) getUserID(c *gin.Context) (uuid.UUID, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return uuid.Nil, registr.ErrInvalidToken
	}

	tokenStr := authHeader[len("Bearer "):]
	claims, err := h.tokenSvc.ValidateToken(c.Request.Context(), tokenStr)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(claims.UserID)
}

func (h *Handler) GetHistory(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, historyModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "UNAUTHORIZED", Message: "invalid token"},
		})
		return
	}

	sessionIDStr := c.Param("id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		h.log.Error("invalid session id", zap.Error(err))
		c.JSON(http.StatusBadRequest, historyModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "INVALID_REQUEST", Message: "invalid session id"},
		})
		return
	}

	session, err := h.sessionRepo.GetByID(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, historyModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "NOT_FOUND", Message: "session not found"},
		})
		return
	}

	if session.CreatorID != userID {
		h.log.Error("unauthorized access to session history")
		c.JSON(http.StatusForbidden, historyModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "FORBIDDEN", Message: "unauthorized"},
		})
		return
	}

	corrections, err := h.service.GetCorrectionsBySessionID(c.Request.Context(), sessionID)
	if err != nil {
		h.log.Error("failed to get corrections", zap.Error(err))
		c.JSON(http.StatusInternalServerError, historyModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "INTERNAL_ERROR", Message: "internal server error"},
		})
		return
	}

	h.log.Info("history fetched successfully", zap.Int("count", len(corrections)))
	c.JSON(http.StatusOK, historyModels.CorrectionsResponse{Corrections: corrections})
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.GET("/history/:id", h.GetHistory)
}