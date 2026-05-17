package route

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	historyService "TZshka/internal/history/service"
	llmModels "TZshka/internal/llm/models"
	llmService "TZshka/internal/llm/service"
	"TZshka/pkg/registr"
)

type Handler struct {
	service        *llmService.Service
	historyService *historyService.Service
	tokenService   *registr.TokenService
	log            *zap.Logger
}

func New(service *llmService.Service, historyService *historyService.Service, tokenService *registr.TokenService, log *zap.Logger) *Handler {
	return &Handler{
		service:        service,
		historyService: historyService,
		tokenService:   tokenService,
		log:            log,
	}
}

func (h *Handler) ProcessText(c *gin.Context) {
	var req llmModels.TextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, llmModels.ErrorResponse{
			Error: struct {
				Type    string `json:"type"`
				Message string `json:"message"`
			}{Type: "validation_error", Message: "invalid request"},
		})
		return
	}

	h.log.Info("processing text request", zap.String("mode", req.Mode), zap.String("standard", req.Standard))

	resp, err := h.service.ProcessText(c.Request.Context(), req.Mode, req.Standard, req.Content)
	if err != nil {
		h.log.Error("failed to process text", zap.Error(err))
		c.JSON(http.StatusInternalServerError, llmModels.ErrorResponse{
			Error: struct {
				Type    string `json:"type"`
				Message string `json:"message"`
			}{Type: "internal_error", Message: "internal server error"},
		})
		return
	}

	if req.SessionID != "" && resp.Success {
		userID := h.getUserIDFromRequest(c)
		sessionID, err := uuid.Parse(req.SessionID)
		if err == nil {
			responseData := make(map[string]interface{})
			if resp.Data != nil {
				responseData = map[string]interface{}{
					"data":    resp.Data,
					"success": resp.Success,
					"code":    resp.Code,
				}
			}
			if userID != uuid.Nil {
				h.historyService.SaveCorrection(c.Request.Context(), sessionID, userID, req.Content, responseData)
			}
		}
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ProcessFile(c *gin.Context) {
	mode := c.PostForm("mode")
	standard := c.PostForm("standard")
	sessionID := c.PostForm("sessionId")

	file, err := c.FormFile("file")
	if err != nil {
		h.log.Error("failed to get file", zap.Error(err))
		c.JSON(http.StatusBadRequest, llmModels.ErrorResponse{
			Error: struct {
				Type    string `json:"type"`
				Message string `json:"message"`
			}{Type: "validation_error", Message: "file is required"},
		})
		return
	}

	opened, err := file.Open()
	if err != nil {
		h.log.Error("failed to open file", zap.Error(err))
		c.JSON(http.StatusInternalServerError, llmModels.ErrorResponse{
			Error: struct {
				Type    string `json:"type"`
				Message string `json:"message"`
			}{Type: "internal_error", Message: "internal server error"},
		})
		return
	}
	defer opened.Close()

	content, err := io.ReadAll(opened)
	if err != nil {
		h.log.Error("failed to read file", zap.Error(err))
		c.JSON(http.StatusInternalServerError, llmModels.ErrorResponse{
			Error: struct {
				Type    string `json:"type"`
				Message string `json:"message"`
			}{Type: "internal_error", Message: "internal server error"},
		})
		return
	}

	h.log.Info("processing file request", zap.String("mode", mode), zap.String("standard", standard), zap.String("filename", file.Filename))

	resp, err := h.service.ProcessText(c.Request.Context(), mode, standard, string(content))
	if err != nil {
		h.log.Error("failed to process file", zap.Error(err))
		c.JSON(http.StatusInternalServerError, llmModels.ErrorResponse{
			Error: struct {
				Type    string `json:"type"`
				Message string `json:"message"`
			}{Type: "internal_error", Message: "internal server error"},
		})
		return
	}

	if sessionID != "" && resp.Success {
		userID := h.getUserIDFromRequest(c)
		sessID, err := uuid.Parse(sessionID)
		if err == nil {
			responseData := make(map[string]interface{})
			if resp.Data != nil {
				responseData = map[string]interface{}{
					"data":    resp.Data,
					"success": resp.Success,
					"code":    resp.Code,
				}
			}
			if userID != uuid.Nil {
				h.historyService.SaveCorrection(c.Request.Context(), sessID, userID, string(content), responseData)
			}
		}
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) getUserIDFromRequest(c *gin.Context) uuid.UUID {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 7 {
		return uuid.Nil
	}

	tokenStr := authHeader[7:]
	claims, err := h.tokenService.ValidateToken(c.Request.Context(), tokenStr)
	if err != nil {
		return uuid.Nil
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil
	}
	return userID
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.POST("/llm/text", h.ProcessText)
	api.POST("/llm/file", h.ProcessFile)
}
