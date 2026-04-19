package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	sessionModels "TZshka/internal/session/models"
	"github.com/gin-gonic/gin"
)

func TestCreateSession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/sessions", func(c *gin.Context) {
		var req sessionModels.CreateSessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, sessionModels.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{Code: "INVALID_REQUEST", Message: "invalid request"},
			})
			return
		}

		if req.Name == "" {
			c.JSON(http.StatusBadRequest, sessionModels.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{Code: "INVALID_REQUEST", Message: "name is required"},
			})
			return
		}

		c.JSON(http.StatusCreated, sessionModels.SessionResponse{
			Session: sessionModels.Session{
				Name: req.Name,
			},
		})
	})

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name: "valid create session",
			body: map[string]string{
				"name": "Test Session",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "missing name",
			body:       map[string]string{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/sessions", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestGetSessions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/sessions", func(c *gin.Context) {
		c.JSON(http.StatusOK, sessionModels.SessionsResponse{
			Sessions: []sessionModels.Session{
				{Name: "Session 1"},
				{Name: "Session 2"},
			},
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/sessions", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestJoinSession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/sessions/join", func(c *gin.Context) {
		var req sessionModels.JoinSessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, sessionModels.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{Code: "INVALID_REQUEST", Message: "invalid request"},
			})
			return
		}

		if req.SessionID == "" {
			c.JSON(http.StatusBadRequest, sessionModels.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{Code: "INVALID_REQUEST", Message: "session id is required"},
			})
			return
		}

		c.JSON(http.StatusOK, sessionModels.SessionResponse{
			Session: sessionModels.Session{
				Name: "Joined Session",
			},
		})
	})

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name: "valid join session",
			body: map[string]string{
				"sessionId": "123e4567-e89b-12d3-a456-426614174000",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing session id",
			body:       map[string]string{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/sessions/join", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}
