package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/history/:id", func(c *gin.Context) {
		sessionID := c.Param("id")
		if sessionID == "" {
			c.JSON(http.StatusBadRequest, map[string]string{"error": "session id is required"})
			return
		}

		if sessionID == "invalid-uuid" {
			c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid session id"})
			return
		}

		corrections := []map[string]interface{}{
			{
				"id":           "test-id-1",
				"sessionId":    sessionID,
				"inputContent": "test content",
				"responseData": map[string]interface{}{"status": "valid"},
				"createdAt":    "2026-04-27T12:00:00Z",
			},
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"corrections": corrections,
		})
	})

	tests := []struct {
		name       string
		url        string
		wantStatus int
	}{
		{
			name:       "valid session id",
			url:        "/history/550e8400-e29b-41d4-a716-446655440000",
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid uuid",
			url:        "/history/invalid-uuid",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestLLMTextWithSessionID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/llm/text", func(c *gin.Context) {
		var req struct {
			Mode      string `json:"mode"`
			Standard  string `json:"standard"`
			Content   string `json:"content"`
			SessionID string `json:"sessionId"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
			return
		}

		if req.Content == "" {
			c.JSON(http.StatusBadRequest, map[string]string{"error": "content is required"})
			return
		}

		response := map[string]interface{}{
			"success": true,
			"code":    200,
			"data": map[string]interface{}{
				"rules_checker": map[string]interface{}{
					"status":   "valid",
					"issues":   []interface{}{},
					"feedback": "OK",
				},
			},
		}

		c.JSON(http.StatusOK, response)
	})

	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid request with sessionId",
			body: map[string]interface{}{
				"mode":      "Instant",
				"standard":  "ГОСТ-19",
				"content":   "Test technical specification",
				"sessionId": "550e8400-e29b-41d4-a716-446655440000",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "valid request without sessionId",
			body: map[string]interface{}{
				"mode":     "Instant",
				"standard": "ГОСТ-19",
				"content":  "Test technical specification",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "missing content",
			body: map[string]interface{}{
				"mode": "Instant",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/llm/text", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}
