package route

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	authModels "TZshka/internal/auth/models"
	authService "TZshka/internal/auth/service"
)

type Handler struct {
	service *authService.Service
	log     *zap.Logger
}

func New(service *authService.Service, log *zap.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req authModels.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, authModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "INVALID_REQUEST", Message: "invalid request"},
		})
		return
	}

	h.log.Info("registering user", zap.String("login", req.Login), zap.String("email", req.Email))

	user, err := h.service.Register(c.Request.Context(), req.Login, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, authService.ErrUserExists) {
			c.JSON(http.StatusBadRequest, authModels.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{Code: "INVALID_REQUEST", Message: "login already exists"},
			})
			return
		}
		h.log.Error("failed to register user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, authModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "INTERNAL_ERROR", Message: "internal server error"},
		})
		return
	}

	h.log.Info("user registered successfully", zap.String("userId", user.ID.String()))
	c.JSON(http.StatusCreated, authModels.UserResponse{User: *user})
}

func (h *Handler) Login(c *gin.Context) {
	var req authModels.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, authModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "INVALID_REQUEST", Message: "invalid request"},
		})
		return
	}

	h.log.Info("logging in user", zap.String("login", req.Login))

	user, token, err := h.service.Login(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, authService.ErrInvalidCreds) {
			c.JSON(http.StatusUnauthorized, authModels.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{Code: "UNAUTHORIZED", Message: "invalid credentials"},
			})
			return
		}
		h.log.Error("failed to login user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, authModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "INTERNAL_ERROR", Message: "internal server error"},
		})
		return
	}

	h.log.Info("user logged in successfully", zap.String("userId", user.ID.String()))
	c.JSON(http.StatusOK, authModels.TokenResponse{Token: token})
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.POST("/register", h.Register)
	api.POST("/login", h.Login)
}
