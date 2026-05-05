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
<<<<<<< HEAD
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

=======
	service  *sessionService.Service
	tokenSvc *registr.TokenService
	log      *zap.Logger
}

func New(service *sessionService.Service, tokenSvc *registr.TokenService, log *zap.Logger) *Handler {
	return &Handler{
		service:  service,
		tokenSvc: tokenSvc,
		log:      log,
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

>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
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
<<<<<<< HEAD
			}{Code: "INVALID_REQUEST", Message: "name is required"},
=======
			}{Code: "INVALID_REQUEST", Message: "invalid request"},
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
		})
		return
	}

<<<<<<< HEAD
	session, err := h.service.Create(c.Request.Context(), req.Name, userID)
=======
	session, err := h.service.CreateSession(c.Request.Context(), req.Name, userID)
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
	if err != nil {
		h.log.Error("failed to create session", zap.Error(err))
		c.JSON(http.StatusInternalServerError, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
<<<<<<< HEAD
			}{Code: "INTERNAL_ERROR", Message: "failed to create session"},
=======
			}{Code: "INTERNAL_ERROR", Message: "internal server error"},
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
		})
		return
	}

<<<<<<< HEAD
=======
	h.log.Info("session created successfully", zap.String("sessionId", session.ID.String()))
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
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

<<<<<<< HEAD
	sessions, err := h.service.GetByUser(c.Request.Context(), userID)
=======
	sessions, err := h.service.GetUserSessions(c.Request.Context(), userID)
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
	if err != nil {
		h.log.Error("failed to get sessions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
<<<<<<< HEAD
			}{Code: "INTERNAL_ERROR", Message: "failed to get sessions"},
=======
			}{Code: "INTERNAL_ERROR", Message: "internal server error"},
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
		})
		return
	}

<<<<<<< HEAD
=======
	h.log.Info("sessions fetched successfully", zap.Int("count", len(sessions)))
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
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
<<<<<<< HEAD
			}{Code: "INVALID_REQUEST", Message: "sessionId is required"},
=======
			}{Code: "INVALID_REQUEST", Message: "invalid request"},
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
		})
		return
	}

	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
<<<<<<< HEAD
		h.log.Error("invalid session id", zap.Error(err))
=======
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
		c.JSON(http.StatusBadRequest, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "INVALID_REQUEST", Message: "invalid session id"},
		})
		return
	}

<<<<<<< HEAD
	if err := h.service.Delete(c.Request.Context(), sessionID, userID); err != nil {
=======
	err = h.service.DeleteSession(c.Request.Context(), sessionID, userID)
	if errors.Is(err, sessionService.ErrSessionNotFound) {
		c.JSON(http.StatusNotFound, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "NOT_FOUND", Message: "session not found"},
		})
		return
	}
	if errors.Is(err, sessionService.ErrUnauthorized) {
		c.JSON(http.StatusForbidden, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "FORBIDDEN", Message: "unauthorized"},
		})
		return
	}
	if err != nil {
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
		h.log.Error("failed to delete session", zap.Error(err))
		c.JSON(http.StatusInternalServerError, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
<<<<<<< HEAD
			}{Code: "INTERNAL_ERROR", Message: "failed to delete session"},
=======
			}{Code: "INTERNAL_ERROR", Message: "internal server error"},
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
		})
		return
	}

<<<<<<< HEAD
	c.JSON(http.StatusOK, gin.H{"message": "session deleted"})
=======
	h.log.Info("session deleted successfully", zap.String("sessionId", sessionID.String()))
	c.JSON(http.StatusOK, sessionModels.SessionResponse{})
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
}

func (h *Handler) JoinSession(c *gin.Context) {
	var req sessionModels.JoinSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
<<<<<<< HEAD
			}{Code: "INVALID_REQUEST", Message: "sessionId is required"},
=======
			}{Code: "INVALID_REQUEST", Message: "invalid request"},
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
		})
		return
	}

	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
<<<<<<< HEAD
		h.log.Error("invalid session id", zap.Error(err))
=======
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
		c.JSON(http.StatusBadRequest, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "INVALID_REQUEST", Message: "invalid session id"},
		})
		return
	}

<<<<<<< HEAD
	session, err := h.service.GetByID(c.Request.Context(), sessionID)
	if err != nil {
		h.log.Error("session not found", zap.Error(err))
=======
	session, err := h.service.JoinSession(c.Request.Context(), sessionID)
	if errors.Is(err, sessionService.ErrSessionNotFound) {
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
		c.JSON(http.StatusNotFound, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "NOT_FOUND", Message: "session not found"},
		})
		return
	}
<<<<<<< HEAD

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
=======
	if err != nil {
		h.log.Error("failed to join session", zap.Error(err))
		c.JSON(http.StatusInternalServerError, sessionModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "INTERNAL_ERROR", Message: "internal server error"},
		})
		return
	}

	h.log.Info("user joined session", zap.String("sessionId", session.ID.String()))
	c.JSON(http.StatusOK, sessionModels.SessionResponse{Session: *session})
>>>>>>> 11ffaa0c6cc4e0527a4a844e359e85b66fdd1f2e
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.POST("/sessions/create", h.CreateSession)
	api.GET("/sessions", h.GetSessions)
	api.DELETE("/sessions", h.DeleteSession)
	api.POST("/sessions/join", h.JoinSession)
}
