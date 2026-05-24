package route

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	sessionModels "TZshka/internal/session/models"
	sessionService "TZshka/internal/session/service"
	"TZshka/pkg/registr"
)

type Handler struct {
	service      *sessionService.Service
	tokenService *registr.TokenService
	log          *zap.Logger
}

func New(service *sessionService.Service, tokenService *registr.TokenService, log *zap.Logger) *Handler {
	return &Handler{
		service:      service,
		tokenService: tokenService,
		log:          log,
	}
}

func (h *Handler) CreateSession(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "UNAUTHORIZED", Message: "invalid token"},
		})
		return
	}

	var req sessionModels.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "INVALID_REQUEST", Message: "name is required"},
		})
		return
	}

	session, err := h.service.Create(c.Request.Context(), req.Name, userID)
	if err != nil {
		h.log.Error("failed to create session", zap.Error(err))
		c.JSON(http.StatusInternalServerError, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "INTERNAL_ERROR", Message: "failed to create session"},
		})
		return
	}

	c.JSON(http.StatusCreated, sessionModels.SessionResponse{Session: *session})
}

func (h *Handler) GetSessions(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "UNAUTHORIZED", Message: "invalid token"},
		})
		return
	}

	sessions, err := h.service.GetByUser(c.Request.Context(), userID)
	if err != nil {
		h.log.Error("failed to get sessions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "INTERNAL_ERROR", Message: "failed to get sessions"},
		})
		return
	}

	c.JSON(http.StatusOK, sessionModels.SessionsResponse{Sessions: sessions})
}

func (h *Handler) DeleteSession(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "UNAUTHORIZED", Message: "invalid token"},
		})
		return
	}

	var req sessionModels.DeleteSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "INVALID_REQUEST", Message: "sessionId is required"},
		})
		return
	}

	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		h.log.Error("invalid session id", zap.Error(err))
		c.JSON(http.StatusBadRequest, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "INVALID_REQUEST", Message: "invalid session id"},
		})
		return
	}

	if err := h.service.Delete(c.Request.Context(), sessionID, userID); err != nil {
		h.log.Error("failed to delete session", zap.Error(err))
		c.JSON(http.StatusInternalServerError, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "INTERNAL_ERROR", Message: "failed to delete session"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session deleted"})
}

func (h *Handler) JoinSession(c *gin.Context) {
	var req sessionModels.JoinSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "INVALID_REQUEST", Message: "sessionId is required"},
		})
		return
	}

	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		h.log.Error("invalid session id", zap.Error(err))
		c.JSON(http.StatusBadRequest, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "INVALID_REQUEST", Message: "invalid session id"},
		})
		return
	}

	session, err := h.service.GetByID(c.Request.Context(), sessionID)
	if err != nil {
		h.log.Error("session not found", zap.Error(err))
		c.JSON(http.StatusNotFound, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "NOT_FOUND", Message: "session not found"},
		})
		return
	}

	c.JSON(http.StatusOK, sessionModels.SessionResponse{Session: *session})
}

func (h *Handler) getUserID(c *gin.Context) (uuid.UUID, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return uuid.Nil, errors.New("no auth header")
	}

	tokenStr := authHeader[7:]
	claims, err := h.tokenService.ValidateToken(c.Request.Context(), tokenStr)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(claims.UserID)
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.POST("/sessions/create", h.CreateSession)
	api.GET("/sessions", h.GetSessions)
	api.DELETE("/sessions", h.DeleteSession)
	api.POST("/sessions/join", h.JoinSession)
}
