package route

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	llmModels "TZshka/internal/llm/models"
	llmService "TZshka/internal/llm/service"
)

type Handler struct {
	service *llmService.Service
	log     *zap.Logger
}

func New(service *llmService.Service, log *zap.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
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

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ProcessFile(c *gin.Context) {
	mode := c.PostForm("mode")
	standard := c.PostForm("standard")

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

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.POST("/llm/text", h.ProcessText)
	api.POST("/llm/file", h.ProcessFile)
}
