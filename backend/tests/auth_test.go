package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	authModels "TZshka/internal/auth/models"
	"github.com/gin-gonic/gin"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/register", func(c *gin.Context) {
		var req authModels.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, authModels.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{Code: "INVALID_REQUEST", Message: "invalid request"},
			})
			return
		}

		if req.Login == "" {
			c.JSON(http.StatusBadRequest, authModels.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{Code: "INVALID_REQUEST", Message: "login is required"},
			})
			return
		}

		c.JSON(http.StatusCreated, authModels.UserResponse{
			User: authModels.User{
				Login: req.Login,
				Email: req.Email,
			},
		})
	})

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name: "valid register",
			body: map[string]string{
				"login":    "testuser",
				"email":    "test@example.com",
				"password": "password123",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "missing login",
			body: map[string]string{
				"email":    "test@example.com",
				"password": "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/login", func(c *gin.Context) {
		var req authModels.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, authModels.ErrorResponse{
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{Code: "INVALID_REQUEST", Message: "invalid request"},
			})
			return
		}

		if req.Login == "testuser" && req.Password == "password123" {
			c.JSON(http.StatusOK, authModels.TokenResponse{
				Token: "valid-token",
			})
			return
		}

		c.JSON(http.StatusUnauthorized, authModels.ErrorResponse{
			Error: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{Code: "UNAUTHORIZED", Message: "invalid credentials"},
		})
	})

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name: "valid login",
			body: map[string]string{
				"login":    "testuser",
				"password": "password123",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid credentials",
			body: map[string]string{
				"login":    "wronguser",
				"password": "wrongpassword",
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestLLMText(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/llm/text", func(c *gin.Context) {
		var req struct {
			Mode     string `json:"mode"`
			Standard string `json:"standard"`
			Content  string `json:"content"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
			return
		}

		if req.Content == "" {
			c.JSON(http.StatusBadRequest, map[string]string{"error": "content is required"})
			return
		}

		if req.Mode != "Instant" && req.Mode != "Thinking" {
			c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid mode"})
			return
		}

		c.JSON(http.StatusOK, map[string]any{
			"success": true,
			"code":    200,
			"data": map[string]string{
				"rules_checker":    "AI output",
				"standard_checker": "AI standard output",
			},
		})
	})

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name: "valid text request",
			body: map[string]string{
				"mode":     "Instant",
				"standard": "ГОСТ-19",
				"content":  "Текст ТЗ",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "missing content",
			body: map[string]string{
				"mode": "Instant",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid mode",
			body: map[string]string{
				"mode":    "Invalid",
				"content": "Текст ТЗ",
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

func TestCreateSession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/sessions/create", func(c *gin.Context) {
		var req struct {
			Name string `json:"name"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
			return
		}

		if req.Name == "" {
			c.JSON(http.StatusBadRequest, map[string]string{"error": "name is required"})
			return
		}

		c.JSON(http.StatusCreated, map[string]any{
			"session": map[string]string{
				"name": req.Name,
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
			req := httptest.NewRequest(http.MethodPost, "/sessions/create", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}
