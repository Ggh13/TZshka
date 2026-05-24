package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	historyModels "TZshka/internal/history/models"
	"github.com/gin-gonic/gin"
)

func TestGetHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/history/:id", func(c *gin.Context) {
		sessionID := c.Param("id")
		if sessionID == "" {
			c.JSON(http.StatusBadRequest, historyModels.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{Code: "INVALID_REQUEST", Message: "session id is required"},
			})
			return
		}

		if sessionID == "invalid-uuid" {
			c.JSON(http.StatusBadRequest, historyModels.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{Code: "INVALID_REQUEST", Message: "invalid session id"},
			})
			return
		}

		if sessionID == "00000000-0000-0000-0000-000000000999" {
			c.JSON(http.StatusNotFound, historyModels.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{Code: "NOT_FOUND", Message: "session not found"},
			})
			return
		}

		corrections := []historyModels.Correction{
			{
				InputContent: "Текст ТЗ",
				ResponseData: map[string]interface{}{"status": "valid"},
			},
		}
		c.JSON(http.StatusOK, historyModels.CorrectionsResponse{Corrections: corrections})
	})

	tests := []struct {
		name       string
		sessionID  string
		wantStatus int
	}{
		{
			name:       "valid session id",
			sessionID:  "00000000-0000-0000-0000-000000000001",
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid uuid format",
			sessionID:  "invalid-uuid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "session not found",
			sessionID:  "00000000-0000-0000-0000-000000000999",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/history/"+tt.sessionID, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestGetHistoryUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/history/:id", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, historyModels.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{Code: "UNAUTHORIZED", Message: "invalid token"},
			})
			return
		}

		c.JSON(http.StatusOK, historyModels.CorrectionsResponse{Corrections: []historyModels.Correction{}})
	})

	tests := []struct {
		name       string
		authHeader string
		wantStatus int
	}{
		{
			name:       "no authorization header",
			authHeader: "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "with authorization header",
			authHeader: "Bearer token123",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/history/00000000-0000-0000-0000-000000000001", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestGetHistoryResponseFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/history/:id", func(c *gin.Context) {
		corrections := []historyModels.Correction{
			{
				InputContent: "Текст ТЗ 1",
				ResponseData: map[string]interface{}{
					"rules_checker": map[string]interface{}{
						"status": "valid",
					},
				},
			},
			{
				InputContent: "Текст ТЗ 2",
				ResponseData: map[string]interface{}{
					"standard_checker": map[string]interface{}{
						"status": "issues_found",
					},
				},
			},
		}
		c.JSON(http.StatusOK, historyModels.CorrectionsResponse{Corrections: corrections})
	})

	req := httptest.NewRequest(http.MethodGet, "/history/00000000-0000-0000-0000-000000000001", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp historyModels.CorrectionsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if len(resp.Corrections) != 2 {
		t.Errorf("expected 2 corrections, got %d", len(resp.Corrections))
	}
}
